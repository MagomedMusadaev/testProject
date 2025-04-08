package config

// Config содержит конфигурационные параметры для сервиса
type Config struct {
	GRPCPort string // порт для gRPC сервера
	DBHost   string // хост базы данных
	DBPort   string // порт базы данных
	DBUser   string // пользователь базы данных
	DBPass   string // пароль базы данных
	DBName   string // имя базы данных
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		GRPCPort: ":50051",
		DBHost:   "postgres", //TODO в реалкейсе нужно всё принимать из .env или yaml файла
		DBPort:   "5432",
		DBUser:   "postgres",
		DBPass:   "postgres",
		DBName:   "telegram_users",
	}
}

// GetDBurl формирует строку подключения к базе данных Psql
func (c *Config) GetDBurl() string {
	return "host=" + c.DBHost + " port=" + c.DBPort +
		" user=" + c.DBUser + " password=" + c.DBPass +
		" dbname=" + c.DBName + " sslmode=disable"
}
