package main

import (
	"context"
	"errors"
	"fmt"
	seed "github.com/hansengotama/disbursement/cmd"
	redishelper "github.com/hansengotama/disbursement/internal/lib/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hansengotama/disbursement/internal/lib/env"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "seed-please" {
		seed.Execute()
		return
	}

	httpHandler := initRoutes()
	initServer(httpHandler)
}

func initServer(httpHandler http.Handler) {
	ctx, cancel := context.WithCancel(context.Background())

	signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	appPort := env.GetAppPort()

	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", appPort),
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  50 * time.Second,
	}

	// run the server in goroutine so that it don't block the main process
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Printf("server started in port %s", server.Addr)

	// Accepts graceful shutdowns when quitting via SIGINT (Ctrl + C)
	// SIGKILL, SIGQUIT or SIGTERM will not be caught and will forcefully shuts the application down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Blocks until we receive graceful shutdown signal
	<-c

	postgres.CloseConnection()
	redishelper.CloseConnection()
	cancel()
	<-ctx.Done()

	log.Print("server is shutting down")
	err := server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}

	log.Print("server has been shut down. goodbye!")
}
