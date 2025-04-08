package domain

import "time"

// User представляет собой модель пользователя в системе
type User struct {
	TelegramID int64     `json:"telegram_id" db:"telegram_id"` // bigint not null unique
	ChannelID  int64     `json:"channel_id" db:"channel_id"`   // bigint not null
	CreatedAt  time.Time `json:"created_at" db:"created_at"`   // timestamp default now()
}
