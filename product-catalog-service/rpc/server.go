package rpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/mosteknoloji/go-grpc-interceptor/panichandler"
	"google.golang.org/grpc"
	"net"
	"product-catalog-service/config"
	"product-catalog-service/controllers"
	"product-catalog-service/integrations/cache/redis"
	"product-catalog-service/integrations/category"
	"product-catalog-service/integrations/database/postgres"
	"product-catalog-service/repositories"
	"product-catalog-service/rpc/interceptors"
	"product-catalog-service/rpc/product"
	"product-catalog-service/services"
	"product-catalog-service/utils"
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

	categoryClient, err := category.NewCategoryClient(config.Configuration)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("categoryClient err conn: %v", err))
	}
	utils.Log("INFO", "Initializing repositories...")
	productRepository := repositories.NewProductRepository(dbConnection)
	redisCache := redis.NewRedisCache("category", redis.GetRedisConnection())

	utils.Log("INFO", "Initializing controllers...")
	productController := controllers.NewProductController(productRepository, redisCache, categoryClient)

	utils.Log("INFO", "Initializing services...")
	productService := services.NewProductCatalogService(productController)

	utils.Log("INFO", "Initializing and registering gRPC services...")
	product.RegisterProductCatalogServiceServer(server, productService)

	utils.Log("INFO", fmt.Sprintf("gRPC server is starting to listen and serve on %s", cfg.Server.RPCAddr))
	go func() {
		if serveErr := server.Serve(rpcListener); serveErr != nil {
			utils.Log("ERROR", fmt.Sprintf("gRPC server stopped with error: %v", serveErr))
			errChan <- serveErr
		}
	}()

	return server, nil
}
