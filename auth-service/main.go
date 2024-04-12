package main

import (
	"auth-service/config"
	"auth-service/utils"
	"fmt"
	"github.com/mosteknoloji/svc"
	"runtime"
)

func main() {
	utils.LogFileInit()
	// Load the application configuration from environment variables.
	cfg, err := config.LoadEnv()
	if err != nil || cfg == nil {
		// Log and exit on configuration load failure.
		utils.Log("FATAL", fmt.Sprintf("Failed to load configuration: %s", err))
		return
	}

	// Initialize the logging file system.

	// Starting the authentication service and logging its HTTP address.
	utils.Log("INFO", fmt.Sprintf("Starting Konzek Senior Backend Developer Assignment Authentication Service @ HTTP API %s", cfg.Server.HTTPAddr))

	// If in debug mode, log the current runtime configuration.
	if config.Configuration.Debug {
		utils.Log("DEBUG", fmt.Sprintf("Runtime config: %s", utils.ToJSON(cfg)))
	}

	// Log the Go runtime version.
	utils.Log("INFO", fmt.Sprintf("Go runtime version: %s", runtime.Version()))

	// Initialize the authentication service.
	srv := new(service)

	// Conditionally start as a Windows Service or a normal service based on configuration.
	if cfg.Server.WinSrv {
		// Log that the service is starting as a Windows Service.
		utils.Log("INFO", "Starting as a Windows Service")
		svc.RunAsService("Konzek Senior Backend Developer Assignment Authentication Service:", srv)
	} else {
		// Log that the service is starting as UNIX Base .
		utils.Log("INFO", "Starting as a UNIX Base Service")
		srv.Start()
	}
}
