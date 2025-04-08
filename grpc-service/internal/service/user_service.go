package service

import (
	"context"
	"errors"
	"f/testProject/api/proto"
	"f/testProject/grpc-service/internal/domain"
	"f/testProject/grpc-service/internal/infrastructure/logger"
	"f/testProject/grpc-service/internal/infrastructure/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

// UserService определяет интерфейс сервиса для работы с пользователями
type UserService interface {
	CheckUserInChannel(ctx context.Context, req *proto.CheckUserRequest) (*proto.CheckUserResponse, error)
}

type userService struct {
	proto.UnimplementedTelegramServiceServer
	userRepo repository.UserRepository
	log      *logger.Logger
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		log:      logger.NewLogger(),
	}
}

func (s *userService) CheckUserInChannel(ctx context.Context, req *proto.CheckUserRequest) (*proto.CheckUserResponse, error) {
	const op = "internal/service/CheckUser"

	s.log.Info("Проверка пользователя", "telegram_id", req.ChannelLink)

	chatID, isMember, err := checkIfUserIsGroup(req.BotToken, req.ChannelLink, req.UserId)
	if err != nil {
		s.log.Error(op, "Ошибка проверки нахождения пользователя в группе", err)
		return nil, err
	}

	// если пользовательне является участником канала
	if !isMember {
		errMsg := "пользователь не является участником канала"
		s.log.Error(
			op,
			errMsg,
			errors.New(errMsg))

		return &proto.CheckUserResponse{
			IsMember: false,
			Error:    errMsg,
		}, nil
	}
	s.log.Info("пользователь подписан на канал")

	user := &domain.User{
		TelegramID: req.UserId,
		ChannelID:  chatID,
		CreatedAt:  time.Now(),
	}

	// Проверка на существование пользователя
	exists, err := s.userRepo.UserExists(req.UserId)
	if err != nil {
		s.log.Error(op, "Ошибка при проверке существования пользователя", err)
		return nil, err
	}

	switch exists {
	case true:
		// Если пользователь уже существует в базе
		s.log.Info("Пользователь уже существует в базе данных")

		return &proto.CheckUserResponse{
			IsMember: true,
		}, nil
	case false:
		// Если пользователя нет в базе, создаем нового
		err = s.userRepo.Create(user)
		if err != nil {
			s.log.Error(op, "Ошибка создания пользователя", err)
			return nil, err
		}
	}

	// Возвращаем успешный ответ
	return &proto.CheckUserResponse{
		IsMember: true,
	}, nil
}

// Функция для проверки, находится ли пользователь в группе
func checkIfUserIsGroup(botToken, channelLink string, userId int64) (int64, bool, error) {
	const op = "internal/service/checkIfUserIsGroup"

	// Создаем клиент Telegram Bot API
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, false, fmt.Errorf("не удалось создать клиент Telegram API: %v", err)
	}

	// Получаем информацию о канале
	chatConfig := tgbotapi.ChatConfig{SuperGroupUsername: channelLink}
	chat, err := bot.GetChat(chatConfig)
	if err != nil {
		return 0, false, fmt.Errorf("не удалось получить информацию о канале: %v", err)
	}

	// Получаем информацию о пользователе в чате
	chatMemberConfig := tgbotapi.ChatConfigWithUser{
		ChatID: chat.ID,
		UserID: int(userId),
	}

	chatMember, err := bot.GetChatMember(chatMemberConfig)
	if err != nil {
		return chat.ID, false, fmt.Errorf("ошибка при вызове Telegram API: %v", err)
	}

	// Если статус члена чата равен "member" или "administrator", значит пользователь в группе
	if chatMember.Status == "member" || chatMember.Status == "administrator" {
		return chat.ID, true, nil
	}

	// Если не член группы, возвращаем false
	return chat.ID, false, nil
}
