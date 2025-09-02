package main

import (
	"effective_mobile/internal/config"
	"effective_mobile/internal/logging.go"
	"effective_mobile/internal/repository"
	"effective_mobile/internal/service"
	"effective_mobile/internal/transport"

	"go.uber.org/zap"
)

func main() {
	config := config.LoadConfig()
	log := logging.InitLogger(config.LogLevel)

	repoManager := repository.NewPostgresReposManager(log)
	if err := repoManager.Connect(config.DSN); err != nil {
		log.Error("connect to db", zap.Error(err))
		return
	}

	handlerManger := transport.NewRouter(
		log,
		service.NewAppServices(repoManager),
	)

	server := transport.NewServer(log, handlerManger.HandlerMap())
	if err := server.Start(config.Address, config.Port); err != nil {
		log.Error("start server", zap.Error(err))
		return
	}
}
