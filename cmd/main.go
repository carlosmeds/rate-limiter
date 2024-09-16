package main

import (
	"fmt"

	"github.com/carlosmeds/rate-limiter/configs"
	"github.com/carlosmeds/rate-limiter/internal/infra/web"
	"github.com/carlosmeds/rate-limiter/internal/infra/web/webserver"
)

func main() {
	configs, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	webserver := webserver.NewWebServer(":" + configs.WebServerPort)
	webOrderHandler := web.NewWebIpHandler()
	webserver.AddHandler("/ip", webOrderHandler.Get)
	fmt.Println("Starting web server on port", webserver.WebServerPort)
	webserver.Start()
}
