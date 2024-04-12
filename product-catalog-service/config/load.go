package config

import (
	"fmt"
	"os"
	"product-catalog-service/utils"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var Configuration *Config

func LoadEnv() (*Config, error) {
	runK8S, err := strconv.ParseBool(os.Getenv("RUN_K8S"))
	if err != nil {
		runK8S = false
	}
	if !runK8S {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}
	debugBoolValue, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("error parsing boolean: %v\n", err))
		return nil, fmt.Errorf("invalid debug variable")
	}
	cfg := Config{
		//ServiceName: "product-catalog-service",
		DEBUG:  debugBoolValue,
		RunK8s: runK8S,
		Server: Server{
			HTTPAddr: os.Getenv("PRODUCT_SERVER_HTTP_ADDR"),
			RPCAddr:  os.Getenv("PRODUCT_SERVER_RPC_ADDR"),
			WinSrv:   runtime.GOOS == "windows",
			Debug:    false,
		},
		Database: Database{
			Host:         os.Getenv("DATABASE_HOST"),
			Port:         os.Getenv("DATABASE_PORT"),
			DatabaseName: os.Getenv("PRODUCT_DATABASE_NAME"),
			User:         os.Getenv("DATABASE_USER"),
			Password:     os.Getenv("DATABASE_PASSWORD"),
		},
		KeyCloak: KeyCloak{
			Realm:                     os.Getenv("KEYCLOAK_REALM"),
			BaseUrl:                   os.Getenv("KEYCLOAK_BASE_URL"),
			RealmKonzekRS256PublicKey: os.Getenv("KEYCLOAK_REALM_KONZEK_RS256_PUBLIC_KEY"),
			KeyCloakClient: KeyCloakClient{
				ClientId:     os.Getenv("KEYCLOAK_PRODUCT_CLIENT_ID"),
				ClientSecret: os.Getenv("KEYCLOAK_PRODUCT_CLIENT_SECRET"),
			},
		},
		CategoryServer: CategoryServer{
			Service:  os.Getenv("CATEGORY_SERVER_SERVICE"),
			HTTPPort: os.Getenv("CATEGORY_SERVER_HTTP_PORT"),
			RPCPort:  os.Getenv("CATEGORY_SERVER_RPC_PORT"),
		},
	}
	Configuration = &cfg
	return &cfg, nil
}
