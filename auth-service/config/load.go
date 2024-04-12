package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"runtime"
	"strconv"
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
		fmt.Printf("error parsing boolean: %v\n", err)
		return nil, fmt.Errorf("invalid debug variable")
	}
	cfg := Config{
		Debug:       debugBoolValue,
		RunK8s:      runK8S,
		ServiceName: "auth-service",
		Server: Server{
			HTTPAddr: os.Getenv("AUTH_SERVER_HTTP_ADDR"),
			WinSrv:   runtime.GOOS == "windows",
			Debug:    false,
		},
		KeyCloak: KeyCloak{
			Realm:                     os.Getenv("KEYCLOAK_REALM"),
			BaseUrl:                   os.Getenv("KEYCLOAK_BASE_URL"),
			RealmMasterRS256PublicKey: os.Getenv("KEYCLOAK_REALM_MASTER_RS256_PUBLIC_KEY"),
			KeyCloakClient: KeyCloakClient{
				ClientId:     os.Getenv("KEYCLOAK_AUTH_CLIENT_ID"),
				ClientSecret: os.Getenv("KEYCLOAK_AUTH_CLIENT_SECRET"),
			},
		},
	}
	Configuration = &cfg
	return &cfg, nil
}
