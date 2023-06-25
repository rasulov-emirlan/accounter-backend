package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/esep-backend/config"
	"github.com/rasulov-emirlan/esep-backend/internal/domains"
	"github.com/rasulov-emirlan/esep-backend/internal/storage/postgresql"
	"github.com/rasulov-emirlan/esep-backend/internal/transport/httprest"
	"github.com/rasulov-emirlan/esep-backend/pkg/logging"
	"github.com/rasulov-emirlan/esep-backend/pkg/validation"
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
	defer log.Sync()

	repo, err := postgresql.NewRepositories(ctx, cfg, log)
	if err != nil {
		log.Fatal("could not init repositories", logging.Error("err", err))
	}
	defer repo.Close()

	commDeps := domains.CommonDependencies{Log: log, Val: validation.GetValidator()}
	authDeps := domains.AuthDependencies{OwnersRepo: repo.Owners(), SecretKey: []byte(cfg.JWTsecret)}
	storesDeps := domains.StoresDependencies{StoresRepo: repo.Stores()}
	doms, err := domains.NewDomainCombiner(commDeps, authDeps, storesDeps)
	if err != nil {
		log.Fatal("could not init domains", logging.Error("err", err))
	}

	srvr := httprest.NewServer(cfg)
	go func() {
		if err := srvr.Start(log, doms); err != nil && err != httprest.ErrServerClosed {
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
