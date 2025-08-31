package main

import (
	"handworks-gateway/internal/config"
	"handworks/common/utils"
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
	go config.StartGQlServer(logger, &wg, stopChan)

	logger.Info("API GATEWAY STARTED SUCCESFULLY")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down account service...")

	close(stopChan)

	// Wait for all routines to clean up
	wg.Wait()
}
