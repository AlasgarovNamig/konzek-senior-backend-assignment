package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/braintree/manners"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"product-categories-service/integrations/cache/redis"
	"product-categories-service/integrations/database/postgres"
	"product-categories-service/repositories"
	"product-categories-service/rest"
	"product-categories-service/rpc"
	"product-categories-service/utils"
	"syscall"
	"time"
)

type service struct {
	rpcServer  *grpc.Server
	restServer *manners.GracefulServer
}

func (s *service) Start() {
	var err error

	utils.Log("INFO", "Initializing database connections...")
	s.InitDatabaseAndCache()
	utils.Log("INFO", "Loading categories to Redis...")
	s.loadCategoriesToRedis()

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

func (s *service) loadCategoriesToRedis() {
	categoryRepository := repositories.NewCategoryRepository(postgres.Connection)
	categories, err := categoryRepository.GetAll()
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to retrieve categories: %v", err))
		panic("Failed to get all categories")
	}
	redisCache := redis.NewRedisCache("category", redis.GetRedisConnection())
	ctx := context.Background()
	err = redisCache.FlushAll(ctx)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to flush Redis: %v", err))
		panic("Failed to flush all in Redis")
	}

	for _, category := range *categories {
		categoryJSON, err := json.Marshal(category)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Error marshaling category: %v", err))
			continue
		}

		err = redisCache.Set(ctx, fmt.Sprintf("%d", category.Id), categoryJSON, 0)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Failed to set category in Redis: %v", err))
			continue
		}
	}
	utils.Log("INFO", "Categories loaded into Redis successfully.")
}
