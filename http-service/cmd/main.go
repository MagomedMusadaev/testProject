package main

import (
	"f/testProject/http-service/internal/config"
	"f/testProject/http-service/internal/infrastructure/logger"
	"f/testProject/http-service/internal/presentation/handlers"
	"f/testProject/http-service/internal/presentation/routes"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	const op = "http-service/cmd/main"

	// Инициализация логгера
	log := logger.NewLogger()

	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Установка соединения с gRPC сервисом
	conn, err := grpc.NewClient(
		cfg.GRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error(op, "Не удалось подключиться к gRPC сервису", err)
		return
	}
	defer conn.Close()

	// Инициализация обработчиков HTTP запросов
	userHandler := handlers.NewUserHandler(conn)

	// Инициализация роутера Gin
	router := gin.Default()

	// Регистрация маршрутов HTTP
	routes.RegisterRoutes(router, userHandler)

	// Запуск HTTP сервера
	log.Info("Запуск HTTP сервера", "порт", cfg.HTTPPort)
	if err = router.Run(cfg.HTTPPort); err != nil {
		log.Error(op, "Не удалось запустить сервер", err)
		return
	}
}
