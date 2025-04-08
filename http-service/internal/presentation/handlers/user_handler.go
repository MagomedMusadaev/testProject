package handlers

import (
	"errors"
	"f/testProject/api/proto"
	"f/testProject/http-service/internal/domain"
	"f/testProject/http-service/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

type UserHandler struct {
	grpcClient *grpc.ClientConn
	log        *logger.Logger
}

func NewUserHandler(grpcClient *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		grpcClient: grpcClient,
		log:        logger.NewLogger(),
	}
}

func (h *UserHandler) CheckUser(c *gin.Context) {
	const op = "internal/presentation/handlers/CheckUser"

	h.log.Info("Обработка запроса пользователя")

	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error(
			op,
			"Ошибка парсинга JSON",
			err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат JSON",
		})
		return
	}

	if user.BotToken == "" || user.ChannelID == "" || user.UserID == 0 {
		h.log.Error(
			op,
			"Не все параметры переданы в запросе",
			errors.New("невалидные параметры запроса"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Параметры bot_token, channel_id и user_id обязательны",
		})
		return
	}

	// Подготовка запроса для gRPC
	req := &proto.CheckUserRequest{
		BotToken:    user.BotToken,
		ChannelLink: user.ChannelID,
		UserId:      int64(user.UserID),
	}

	// Создаём gRPC клиент
	client := proto.NewTelegramServiceClient(h.grpcClient)

	// Вызываем gRPC метод и обрабатываем ответ
	resp, err := client.CheckUserInChannel(c, req)
	if err != nil {
		h.log.Error(op, "Ошибка при вызове gRPC сервиса: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при обработке запроса",
		})
		return
	}

	// Логируем и обрабатываем результат
	switch resp.IsMember {
	case true:
		h.log.Info("Пользователь состоит в канале")
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Пользователь состоит в канале",
		})
	case false:
		h.log.Info("Пользователь не состоит в канале")
		c.JSON(http.StatusOK, gin.H{
			"status":  "failure",
			"message": "Пользователь не состоит в канале",
		})
	}
	h.log.Info("Проверка пользователя завершена")
}
