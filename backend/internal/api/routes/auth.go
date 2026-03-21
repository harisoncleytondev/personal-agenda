package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harisoncleytondev/personal-agenda/internal/api/routes/dto"
	"github.com/harisoncleytondev/personal-agenda/internal/service"
)

type LoginHandler struct {
	authService *service.AuthService
}

func NewAuthHandle (authService *service.AuthService) *LoginHandler {
	return &LoginHandler{authService: authService}
}

func (l *LoginHandler) Login(c * gin.Context) {
	var body *dto.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := l.authService.AuthenticationUserLogin(c.Request.Context(), body)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", access, 24*60*60, "/", "", false, true)
	c.SetCookie("refresh_token", refresh, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login realizado com sucesso"})
}

func (l *LoginHandler) Register(c * gin.Context) {
	var body *dto.RegisterRequest

	if err :=c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := l.authService.AuthenticationUserRegister(c.Request.Context(), body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Conta criada com sucesso."})
}

func (l *LoginHandler) Refresh(c *gin.Context) {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token não encontrado"})
		return
	}

	access, refresh, err := l.authService.RefreshToken(c.Request.Context(), refreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", access, 24*60*60, "/", "", false, true)
	c.SetCookie("refresh_token", refresh, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Tokens atualizados"})
}