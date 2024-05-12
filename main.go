package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/rexcfnghk/pricing-store/application"
	"github.com/rexcfnghk/pricing-store/config"
)

func main() {
	appConfig := readAppConfig()

	app := application.New(appConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}

func readAppConfig() *config.AppConfig {
	confFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("failed to open config file", err)
	}
	defer func(confFile *os.File) {
		err := confFile.Close()
		if err != nil {
			fmt.Println("failed to close config file", err)
		}
	}(confFile)
	conf, err := io.ReadAll(confFile)
	if err != nil {
		fmt.Println("failed to read config file", err)
	}

	appConfig := &config.AppConfig{}
	err = json.Unmarshal(conf, &appConfig)
	if err != nil {
		fmt.Println("failed to parse config file into JSON", err)
	}

	return appConfig
}
