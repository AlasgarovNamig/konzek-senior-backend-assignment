package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/braintree/manners"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"product-catalog-service/api"
	"product-catalog-service/config"
	"product-catalog-service/rest/middlewares"
	"product-catalog-service/rpc/product"
	"product-catalog-service/utils"
)

func customError(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	utils.Log("INFO", fmt.Sprintf("Handling error for request %v: %v", r.URL.Path, err))
	customErr := api.GetGRPCErrorResponse(err)
	utils.Log("INFO", fmt.Sprintf("Mapped gRPC error to HTTP status %d for request %v", customErr.Status, r.URL.Path))

	w.WriteHeader(customErr.Status)
	if encodeErr := json.NewEncoder(w).Encode(customErr); encodeErr != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to encode custom error response for request %v: %v", r.URL.Path, encodeErr))
	}
}

func NewRestServer(errChan chan error) (*manners.GracefulServer, error) {
	utils.Log("INFO", "Initializing REST server...")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux(runtime.WithErrorHandler(customError))
	utils.Log("INFO", "Custom error handler set for REST server.")

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cfg := config.Configuration
	server := manners.NewWithServer(&http.Server{
		Addr: cfg.Server.HTTPAddr,
		Handler: middlewares.RequestLatencyTrackingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				httpMux.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		})),
	})
	utils.Log("INFO", fmt.Sprintf("REST server configured to listen on %s", cfg.Server.HTTPAddr))

	err := product.RegisterProductCatalogServiceHandlerFromEndpoint(ctx, mux, cfg.Server.RPCAddr, opts)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to register gRPC service handlers with the REST server: %v", err))
		defer cancel()
		return nil, err
	}
	utils.Log("INFO", "gRPC service handlers registered with the REST server successfully.")
	go func() {
		utils.Log("INFO", fmt.Sprintf("REST server is starting to listen and serve on %s", cfg.Server.HTTPAddr))
		defer cancel()
		if err := server.ListenAndServe(); err != nil {
			utils.Log("ERROR", fmt.Sprintf("REST server stopped with error: %v", err))
			errChan <- err
		}
	}()
	return server, nil
}
