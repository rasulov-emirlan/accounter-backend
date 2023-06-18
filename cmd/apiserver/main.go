package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/esep-backend/config"
	"github.com/rasulov-emirlan/esep-backend/internal/transport/httprest"
	"github.com/rasulov-emirlan/esep-backend/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logging.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	srvr := httprest.NewServer(cfg.Server.Port)
	go func() {
		if err := srvr.Start(log); err != nil && err != httprest.ErrServerClosed {
			log.Fatal("server start", logging.Error("err", err))
		}
	}()

	log.Info("server started", logging.String("port", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("graceful shutdown")

	if err := srvr.Stop(ctx); err != nil {
		log.Fatal("server stop", logging.Error("err", err))
	}
}
