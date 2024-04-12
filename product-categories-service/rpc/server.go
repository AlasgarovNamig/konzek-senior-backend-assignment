package rpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/mosteknoloji/go-grpc-interceptor/panichandler"
	"google.golang.org/grpc"
	"net"
	"product-categories-service/config"
	"product-categories-service/controllers"
	"product-categories-service/integrations/cache/redis"
	"product-categories-service/integrations/database/postgres"
	"product-categories-service/repositories"
	"product-categories-service/rpc/category"
	"product-categories-service/rpc/interceptors"
	"product-categories-service/services"
	"product-categories-service/utils"
)

func NewRPCAPIServer(errChan chan error) (*grpc.Server, error) {
	utils.Log("INFO", "Initializing gRPC server...")

	cfg := config.Configuration

	utils.Log("INFO", fmt.Sprintf("Attempting to listen on RPC address: %s", cfg.Server.RPCAddr))
	rpcListener, err := net.Listen("tcp", cfg.Server.RPCAddr)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to listen on %s: %v", cfg.Server.RPCAddr, err))
		return nil, err
	}

	utils.Log("INFO", "Installing panic handler...")
	panichandler.InstallPanicHandler(interceptors.LogPanicStackMultiLine)
	utils.Log("INFO", "Configuring server options...")
	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			panichandler.UnaryServerInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			interceptors.UnaryLogHandler,
			interceptors.UnaryAuthHandler,
		)),
	}

	utils.Log("INFO", "Creating the gRPC server with configured options...")
	server := grpc.NewServer(options...)

	utils.Log("INFO", "Establishing database connection...")
	dbConnection := postgres.Connection

	utils.Log("INFO", "Initializing repositories...")
	categoryRepository := repositories.NewCategoryRepository(dbConnection)
	redisCache := redis.NewRedisCache("category", redis.GetRedisConnection())

	utils.Log("INFO", "Initializing controllers...")
	categoryController := controllers.NewCategoryController(categoryRepository, redisCache)

	utils.Log("INFO", "Initializing services...")
	categoryService := services.NewProductCategoryService(categoryController)

	utils.Log("INFO", "Initializing and registering gRPC services...")
	category.RegisterProductCategoryServiceServer(server, categoryService)

	utils.Log("INFO", fmt.Sprintf("gRPC server is starting to listen and serve on %s", cfg.Server.RPCAddr))
	go func() {
		if serveErr := server.Serve(rpcListener); serveErr != nil {
			utils.Log("ERROR", fmt.Sprintf("gRPC server stopped with error: %v", serveErr))
			errChan <- serveErr
		}
	}()

	return server, nil
}
