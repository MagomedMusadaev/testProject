package routes

import (
	"f/testProject/http-service/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(router *gin.Engine, handler *handlers.UserHandler) {
	RegisterUserRoutes(router, handler)
}

// RegisterUserRoutes регистрирует маршруты для пользовательских эндпоинтов
func RegisterUserRoutes(router *gin.Engine, handler *handlers.UserHandler) {
	router.POST("/check-user", handler.CheckUser)
}
