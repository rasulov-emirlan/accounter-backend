package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/accounter-backend/config"
	"github.com/rasulov-emirlan/accounter-backend/internal/domains"
	"github.com/rasulov-emirlan/accounter-backend/internal/storage/postgresql"
	"github.com/rasulov-emirlan/accounter-backend/internal/transport/httprest"
	"github.com/rasulov-emirlan/accounter-backend/pkg/logging"
	"github.com/rasulov-emirlan/accounter-backend/pkg/shutdown"
	"github.com/rasulov-emirlan/accounter-backend/pkg/telemetry"
	"github.com/rasulov-emirlan/accounter-backend/pkg/validation"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cleaner := shutdown.NewScheduler()

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
	cleaner.Add(repo.Close)
	log.Info("repositories initialized")

	commDeps := domains.CommonDependencies{Log: log, Val: validation.GetValidator()}
	authDeps := domains.AuthDependencies{OwnersRepo: repo.Owners(), SecretKey: []byte(cfg.JWTsecret)}
	storesDeps := domains.StoresDependencies{StoresRepo: repo.Stores()}
	categoriesDeps := domains.CategoriesDependencies{CategoriesRepo: repo.Categories()}
	doms, err := domains.NewDomainCombiner(commDeps, authDeps, storesDeps, categoriesDeps)
	if err != nil {
		log.Fatal("could not init domains", logging.Error("err", err))
	}
	log.Info("domains initialized")

	srvr := httprest.NewServer(cfg)
	cleaner.Add(srvr.Stop)
	go func() {
		if err := srvr.Start(log, doms); err != nil && err != httprest.ErrServerClosed {
			log.Fatal("server start", logging.Error("err", err))
		}
	}()
	log.Info("server started", logging.String("port", cfg.Server.Port))

	if err := telemetry.StartJaegerTraceProvider(cfg, &cleaner); err != nil {
		log.Fatal("telemetry start", logging.Error("err", err))
	}
	log.Info("telemetry started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("graceful shutdown")

	if err := cleaner.Close(ctx); err != nil {
		log.Fatal("cleaner close", logging.Error("err", err))
	}
}
