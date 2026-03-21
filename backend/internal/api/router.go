package api

import (
	"github.com/gin-gonic/gin"
	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/api/routes"
	"github.com/harisoncleytondev/personal-agenda/internal/repository"
	"github.com/harisoncleytondev/personal-agenda/internal/service"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := routes.NewAuthHandle(authService)

	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/register", authHandler.Register)
	r.GET("/auth/refresh", authHandler.Refresh)

	return r
}