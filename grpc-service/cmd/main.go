package main

import (
	"database/sql"
	"net"

	"f/testProject/api/proto"
	"f/testProject/grpc-service/internal/config"
	"f/testProject/grpc-service/internal/infrastructure/logger"
	"f/testProject/grpc-service/internal/infrastructure/repository"
	"f/testProject/grpc-service/internal/service"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	const op = "grpc-service/cmd/main"

	// Инициализация логгера
	log := logger.NewLogger()

	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Установка соединения с базой данных Psql
	db, err := sql.Open("postgres", cfg.GetDBurl())
	if err != nil {
		log.Error(op, "Не удалось подключиться к базе данных", err)
		return
	}
	defer db.Close()

	// Инициализация репозитория для работы с пользователями
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервиса бизнес-логики
	userService := service.NewUserService(userRepo)

	// Инициализация TCP-слушателя для gRPC сервера
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Error(op, "Не удалось начать прослушивание порта", err)
		return
	}

	// Создание и настройка gRPC сервера
	s := grpc.NewServer()
	proto.RegisterTelegramServiceServer(s, userService)

	// Запуск gRPC сервера
	log.Info("Запуск gRPC сервера", "порт", cfg.GRPCPort)
	if err = s.Serve(lis); err != nil {
		log.Error(op, "Не удалось запустить сервер", err)
		return
	}
}
