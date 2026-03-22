package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/repository"
)

func AuthMiddleware(userRepo *repository.UserRepository, tokenTypeExpected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Método de assinatura inválido")
			}
			return []byte(config.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		tType, ok := claims["type"].(string)
		if !ok || tType != tokenTypeExpected {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		user, err := userRepo.FindById(c.Request.Context(), sub)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("user", user)

		c.Next()
	}
}