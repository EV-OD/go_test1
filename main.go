package main

import (
	"fmt"
	"myapp/config"
	"myapp/router"
	"os"
)

func main() {

	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	e := router.NewRouter(cfg)
	address := fmt.Sprintf(":%v", cfg.Server.Port)
	if err := e.Start(address); err != nil {
		e.Logger.Error("server failed to start", "error", err)
		os.Exit(1)
	} else {
		e.Logger.Info("server started successfully", "port", cfg.Server.Port)
	}

}
