package main

import (
	"cachepractice/cache"
	"cachepractice/envconfig"
	"cachepractice/logger"
	"cachepractice/repository"
	"cachepractice/server"
	"fmt"
	"log"

	"go.uber.org/zap"
)

var (
	env       envconfig.Env
	db        repository.Repository
	zaplogger *zap.Logger
)

func main() {
	var err error
	zaplogger, err = logger.Init()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}

	env, err = envconfig.Process()
	if err != nil {
		log.Fatalf("failed to process env: %s", err)
	}

	db, err = repository.InitFakeDB()
	if err != nil {
		log.Fatalf("failed to connect db: %s", err)
	}

	cachedb := cache.InitCacheLayer(db, zaplogger)

	r := server.NewRouter(cachedb, zaplogger)
	r.Run(fmt.Sprintf(":%d", env.AppPort))
}
