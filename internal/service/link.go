package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Wild-sergunys/shrtic/internal/model"
	"github.com/Wild-sergunys/shrtic/internal/repository"
)

const (
	shortCodeLength = 7
	base62Chars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	cacheTTL        = 24 * time.Hour
)

type LinkService struct {
	linkRepo *repository.LinkRepository
	redis    *redis.Client
}

func NewLinkService(linkRepo *repository.LinkRepository, redis *redis.Client) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
		redis:    redis,
	}
}

func generateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", fmt.Errorf("ошибка генерации случайного числа: %w", err)
		}
		code[i] = base62Chars[n.Int64()]
	}
	return string(code), nil
}

func (s *LinkService) CreateShortLink(ctx context.Context, longURL string, userID *int64) (*model.Link, error) {
	longURL = strings.TrimSpace(longURL)
	if longURL == "" {
		return nil, fmt.Errorf("URL не может быть пустым")
	}

	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "https://" + longURL
	}

	code, err := generateShortCode()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации кода: %w", err)
	}

	link := &model.Link{
		UserID:   userID,
		ShortURL: "/" + code,
		LongURL:  longURL,
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("ошибка создания ссылки: %w", err)
	}

	if err := s.redis.Set(ctx, "/"+code, longURL, cacheTTL).Err(); err != nil {
		return nil, fmt.Errorf("ошибка кэширования: %w", err)
	}

	return link, nil
}

func (s *LinkService) GetLongURL(ctx context.Context, code string) (string, error) {
	longURL, err := s.redis.Get(ctx, code).Result()
	if err == nil {
		return longURL, nil
	}

	link, err := s.linkRepo.FindByShortCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска ссылки: %w", err)
	}
	if link == nil {
		return "", fmt.Errorf("ссылка не найдена")
	}

	s.redis.Set(ctx, code, link.LongURL, cacheTTL)

	return link.LongURL, nil
}

func (s *LinkService) GetLinks(ctx context.Context, userID int64, search string) ([]model.Link, error) {
	return s.linkRepo.FindByUserID(ctx, userID, search)
}

func (s *LinkService) DeleteLink(ctx context.Context, userID int64, linkID int64) error {
	link, err := s.linkRepo.FindByID(ctx, linkID)
	if err != nil {
		return fmt.Errorf("ошибка поиска ссылки: %w", err)
	}
	if link == nil {
		return fmt.Errorf("ссылка не найдена")
	}
	if link.UserID == nil || *link.UserID != userID {
		return fmt.Errorf("нет прав на удаление этой ссылки")
	}

	if err := s.linkRepo.Delete(ctx, linkID); err != nil {
		return fmt.Errorf("ошибка удаления ссылки: %w", err)
	}

	code := link.ShortURL
	s.redis.Del(ctx, code)

	return nil
}

func (s *LinkService) RecordClick(ctx context.Context, linkID int64) error {
	return s.linkRepo.IncrementClicks(ctx, linkID)
}
