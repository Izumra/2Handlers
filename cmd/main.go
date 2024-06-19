package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Izumra/2Handlers/docs"
	"github.com/Izumra/2Handlers/internal/server"
	"github.com/Izumra/2Handlers/internal/service/auth"
	"github.com/Izumra/2Handlers/internal/service/news"
	"github.com/Izumra/2Handlers/internal/storage/postgresql"
	"github.com/Izumra/2Handlers/utils/config"
	"github.com/Izumra/2Handlers/utils/logger"
)

// @title 2 Обработчика
// @version        1.0
// @description    2 обработчика.

// @BasePath /

// @securityDefinitions.apikey Authorization
// @in header
// @name authorization
func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Ошибка при генерации конфигурации", slog.Any("ошибка", err))
		os.Exit(1)
	}
	log.Info("Считанная конфигурация", slog.Any("конфиг", cfg))

	db, err := postgresql.New(context.Background(), cfg.Db.ConnString)
	if err != nil {
		log.Error("Ошибка при подключении к базе данных", slog.Any("ошибка", err))
		os.Exit(1)
	}

	newsService := news.New(log, cfg.Token, db)
	authService := auth.New(log, db, cfg.Token)
	server := server.New(newsService, authService)
	server.RegHandlers()

	chanServ := make(chan error)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go server.Start(addr, chanServ)

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sign := <-exitChan:
		err := server.Stop()
		if err != nil {
			log.Error("Ошибка при остановке сервера перед завершением работы программы", slog.Any("ошибка", err))
			os.Exit(1)
		}

		log.Info("Программа завершила работу, сервер правильно остановлен", slog.Any("сигнал", sign))
	case err := <-chanServ:
		log.Error("Ошибка при запуске сервера", slog.Any("ошибка", err))
	}
}
