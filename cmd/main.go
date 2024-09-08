package main

import (
	"fmt"

	"github.com/carlosmeds/rate-limiter/internal/infra/web"
	"github.com/carlosmeds/rate-limiter/internal/infra/web/webserver"
)

func main() {
	fmt.Println("Hello, World!")

	webserver := webserver.NewWebServer(":8080")
	webOrderHandler := web.NewWebIpHandler()
	webserver.AddHandler("/ip", webOrderHandler.Get)
	fmt.Println("Starting web server on port", webserver.WebServerPort)
	webserver.Start()
}
