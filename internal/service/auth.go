package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Wild-sergunys/shrtic/internal/model"
	"github.com/Wild-sergunys/shrtic/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
	jwtTTL   time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, jwtKey string, jwtTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
		jwtTTL:   jwtTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (*model.User, error) {
	if len(password) < 6 {
		return nil, fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	existing, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке пользователя: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("пользователь с таким логином уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хэширования пароля: %w", err)
	}

	user := &model.User{
		Login:    login,
		Password: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(s.jwtTTL).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("ошибка создания токена: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

var (
	ErrInvalidCredentials = &AuthError{Message: "Неверный логин или пароль"}
)

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
