package api

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/api/routes"
	"github.com/harisoncleytondev/personal-agenda/internal/jobs"
	"github.com/harisoncleytondev/personal-agenda/internal/middleware"
	"github.com/harisoncleytondev/personal-agenda/internal/repository"
	"github.com/harisoncleytondev/personal-agenda/internal/service"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userRepo := repository.NewUserRepository(config.DB)

	authService := service.NewAuthService(userRepo)
	authHandler := routes.NewAuthHandle(authService)

	appointmentRepo := repository.NewAppointmentRepository(config.DB)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	appointmentHandler := routes.NewAppointmentHandle(appointmentService)

	jobs.Start(appointmentService)
	log.Println("Cron Jobs iniciados com sucesso.")

	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/register", authHandler.Register)
	r.GET("/auth/refresh", authHandler.Refresh)

	authGroup := r.Group("/logged/appointment")
	authGroup.Use(middleware.AuthMiddleware(userRepo, "access"))
	{
		authGroup.GET("/", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
		authGroup.POST("/create", appointmentHandler.Create)
		authGroup.GET("/getall", appointmentHandler.GetAll)
		authGroup.GET("/get/:id", appointmentHandler.GetByID)
		authGroup.PUT("/update/:id", appointmentHandler.Update)
		authGroup.DELETE("/delete/:id", appointmentHandler.Delete)
	}

	return r
}