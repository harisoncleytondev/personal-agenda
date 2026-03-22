package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/api/routes/dto"
	"github.com/harisoncleytondev/personal-agenda/internal/model"
	"github.com/harisoncleytondev/personal-agenda/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService (repo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func generateTokens(userID string) (accessToken string, refreshToken string, err error) {
    now := time.Now()

    accessClaims := jwt.MapClaims{
        "sub":  userID,
        "exp":  now.Add(24 * time.Hour).Unix(),
        "type": "access",
    }
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessToken, err = at.SignedString([]byte(config.GetJWTSecret())) 
    if err != nil {
        return "", "", err
    }

    refreshClaims := jwt.MapClaims{
        "sub":  userID,
        "exp":  now.Add(30 * 24 * time.Hour).Unix(),
        "type": "refresh",
    }
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshToken, err = rt.SignedString([]byte(config.GetJWTSecret()))
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

func (s *AuthService) AuthenticationUserLogin(ctx context.Context, user *dto.LoginRequest) (accessToken string, refreshToken string, err error) {
	userFind, err := s.userRepo.FindByEmail(ctx, user.Email) 

	if err != nil {
		fmt.Println("Erro ao buscar usuário:", err)
		return "", "", errors.New("Usuário não encontrado.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFind.PasswordHash), []byte(user.Password))

	if err != nil {
		return "", "", errors.New("Senha ou email incorretos.")
	}

	return generateTokens(userFind.ID)
}

func (s *AuthService) AuthenticationUserRegister(ctx context.Context, user *dto.RegisterRequest) error {
	if user.Password != user.PasswordCO {
		return errors.New("As senhas não coincidem")
	}

	_, err := s.userRepo.FindByEmail(ctx, user.Email)

	if err == nil {
		return errors.New("e-mail já está em uso")
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	newUser := &model.UserCreate{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: string(passwordHash), 
	}

	err = s.userRepo.UserCreate(ctx, newUser)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Método de assinatura inválido")
		}
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("Refresh token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", "", errors.New("Refresh token inválido")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("Refresh token inválido")
	}

	user, err := s.userRepo.FindByEmail(ctx, sub)
	if err != nil || user == nil {
		return "", "", errors.New("Usuário não encontrado")
	}

	return generateTokens(user.ID)
}