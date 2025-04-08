package repository

import (
	"database/sql"
	"errors"
	"f/testProject/grpc-service/internal/domain"
	"f/testProject/grpc-service/internal/infrastructure/logger"
	"time"
)

// UserRepository определяет интерфейс для работы с хранилищем пользователей
type UserRepository interface {
	Create(user *domain.User) error
	UserExists(telegramID int64) (bool, error)
}

type userRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db:  db,
		log: logger.NewLogger(),
	}
}

func (r *userRepository) Create(user *domain.User) error {
	const op = "internal/infrastructure/repository/Create"

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	query := `INSERT INTO users (telegram_id, channel_id, created_at) VALUES ($1, $2, $3) RETURNING id`

	r.log.Info("Создание нового пользователя", "telegram_id", user.TelegramID)

	var userId int

	err := r.db.QueryRow(
		query,
		user.TelegramID,
		user.ChannelID,
		user.CreatedAt).Scan(&userId)
	if err != nil {
		r.log.Error(op, "Ошибка создания пользователя", err)
		return err
	}

	r.log.Info("Пользователь успешно создан", "user_id", userId, "telegram_id", user.TelegramID)
	return nil
}

func (r *userRepository) UserExists(telegramID int64) (bool, error) {
	const op = "internal/infrastructure/repository/UserExists"

	// Запрос для проверки существования пользователя с данным telegram_id
	var existingUserID int
	query := `SELECT id FROM users WHERE telegram_id = $1`
	err := r.db.QueryRow(query, telegramID).Scan(&existingUserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Пользователь не найден, возвращаем false
			return false, nil
		}
		r.log.Error(op, "Ошибка проверки пользователя", err)
		return false, err
	}

	// Если пользователь найден, возвращаем true
	return true, nil
}
