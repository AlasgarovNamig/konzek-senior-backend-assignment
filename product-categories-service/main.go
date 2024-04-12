package main

import (
	"flag"
	"fmt"
	"github.com/mosteknoloji/svc"
	"product-categories-service/config"
	"product-categories-service/utils"
	"runtime"
)

func main() {
	flag.Parse()
	utils.LogFileInit()
	utils.Log("INFO", "Loading configuration settings from environment variables.")
	cfg, err := config.LoadEnv()
	if err != nil || cfg == nil {
		utils.Log("ERROR", fmt.Sprintf("[FATAL] Failed to load configuration: %s", err))
		return
	}

	utils.Log("INFO", fmt.Sprintf("Starting Konzek Senior Backend Developer Assignment Authentication Service @ HTTP API %s", cfg.Server.HTTPAddr))
	if cfg.DEBUG {
		utils.Log("DEBUG", fmt.Sprintf("Runtime config: %s", utils.ToJSON(cfg)))
	}
	utils.Log("INFO", fmt.Sprintf("Go runtime version: %s", runtime.Version()))

	srv := new(service)
	if cfg.Server.WinSrv {
		utils.Log("INFO", "Running as a Windows Service")
		svc.RunAsService("Product Categories Service", srv)
	} else {
		utils.Log("INFO", "Starting the service as a UNIX Base.")
		srv.Start()
	}
}
