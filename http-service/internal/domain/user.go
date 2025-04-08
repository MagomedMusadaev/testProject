package domain

// User представляет собой модель пользователя
type User struct {
	BotToken  string `json:"bot_token"`
	ChannelID string `json:"channel_link"`
	UserID    int    `json:"user_id"`
}
