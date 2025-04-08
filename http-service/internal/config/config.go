package config

// Config содержит конфигурационные параметры HTTP-сервиса
type Config struct {
	HTTPPort     string // HTTPPort определяет порт для HTTP-сервера
	GRPCEndpoint string // GRPCEndpoint определяет адрес подключения к gRPC-сервису
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		HTTPPort:     ":8085", //TODO: по хорошему эти данные должны быть в .env файле!
		GRPCEndpoint: "grpc-service:50051",
	}
}
