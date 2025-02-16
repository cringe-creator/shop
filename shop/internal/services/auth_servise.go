package services

import (
	"errors"
	"os"
	"time"

	"shop/internal/models"
	"shop/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type customTokenClaims struct {
	jwt.RegisteredClaims
	UserID int   `json:"user_id"`
	exp    int64 `json:"exp"`
}

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func (s *AuthService) Register(username, password string) (*models.User, string, error) {
	existingUser, _ := s.UserRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, "", errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.UserRepo.CreateUser(username, string(hashedPassword))
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	return s.GenerateToken(user.ID)
}

func (s *AuthService) GenerateToken(userID int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // Лучше вынести в переменные окружения
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
