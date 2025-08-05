package main

import (
	"handworks/services/inventory/handlers"
	"handworks/services/inventory/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load()
	if err != nil {
		logger.Info("No .env file found (this is okay in production)")
	}
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	wg.Add(1)
	go handlers.StartGQlServer(logger, &wg, stopChan)

	logger.Info("Inventory service started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down account service...")

	close(stopChan)

	// Wait for all routines to clean up
	wg.Wait()
}
