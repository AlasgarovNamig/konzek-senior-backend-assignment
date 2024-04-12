package main

import (
	"context"
	"fmt"
	"github.com/braintree/manners"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"product-catalog-service/integrations/cache/redis"
	"product-catalog-service/integrations/database/postgres"
	"product-catalog-service/rest"
	"product-catalog-service/rpc"
	"product-catalog-service/utils"
	"syscall"
	"time"
)

type service struct {
	rpcServer  *grpc.Server
	restServer *manners.GracefulServer
}

func (s *service) Start() {
	var err error

	utils.Log("INFO", "Initializing database and cache connections...")
	s.InitDatabaseAndCache()

	errRPCChan := make(chan error, 1)
	utils.Log("INFO", "Starting RPC server...")
	s.rpcServer, err = rpc.NewRPCAPIServer(errRPCChan)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Could not start RPC server: %s", err))
		return
	}

	errRestChan := make(chan error, 1)
	utils.Log("INFO", "Starting REST server...")
	s.restServer, err = rest.NewRestServer(errRestChan)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Could not start REST server: %+v", err))
		panic("Could not initialize REST server")
	}

	utils.Log("INFO", "Service started successfully. Waiting for termination signal...")
	s.waitSignal(errRPCChan)
}

func (s *service) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	utils.Log("INFO", "Shutting down REST server...")
	err := s.restServer.Shutdown(ctx)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Error shutting down REST server: %+v", err))
	}

	utils.Log("INFO", "Stopping RPC server...")
	s.rpcServer.GracefulStop()

	utils.Log("INFO", "Closing database and Redis connections...")
	postgres.CloseDatabaseConnection()
	redis.CloseRedisConnection()
	utils.LogFile.Close()
	utils.Log("INFO", "Service stopped.")
}

func (s *service) InitDatabaseAndCache() {
	postgres.SetupDatabaseConnection()
	utils.Log("INFO", "Database connection established.")

	redis.SetupRedisConnection()
	utils.Log("INFO", "Redis connection established.")
}

func (s *service) waitSignal(errRPCChan chan error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case err := <-errRPCChan:
			if err != nil {
				utils.Log("ERROR", fmt.Sprintf("RPC server error: %+v", err))
			}
		case <-signalChan:
			utils.Log("INFO", "Signal received. Initiating shutdown...")
			s.Stop()
			return
		}
	}
}
