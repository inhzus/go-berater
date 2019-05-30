package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/inhzus/go-berater/config"
	"github.com/inhzus/go-berater/models"
	"github.com/inhzus/go-berater/routes"
	"os"
	"os/signal"
)

func setup() {
	config.Setup()
	models.Setup()
}

func teardown() {
	models.Teardown()
}

func main() {
	setup()
	engine := gin.Default()
	routes.ApplyRoutes(engine)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		fmt.Println("\nReceived an interrupt, teardown")
		teardown()
		os.Exit(0)
	}()
	_ = engine.Run(":5000")
}
