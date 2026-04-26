package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
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
		ShortURL: "/r/" + code,
		LongURL:  longURL,
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("ошибка создания ссылки: %w", err)
	}

	if err := s.redis.Set(ctx, "/r/"+code, longURL, cacheTTL).Err(); err != nil {
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

func (s *LinkService) GetLinkByCode(ctx context.Context, code string) (*model.Link, error) {
	link, err := s.linkRepo.FindByShortCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска ссылки: %w", err)
	}
	if link == nil {
		return nil, fmt.Errorf("ссылка не найдена")
	}
	return link, nil
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

	s.redis.Del(ctx, link.ShortURL)

	return nil
}

func (s *LinkService) RecordClick(ctx context.Context, linkID int64) error {
	return s.linkRepo.IncrementClicks(ctx, linkID)
}

func (s *LinkService) SaveClickStats(ctx context.Context, linkID int64, userAgent, referer, clientIP string) error {
	browser, device := parseUserAgent(userAgent)
	country := getCountry(clientIP)
	source := parseReferer(referer)

	_, err := s.linkRepo.SaveClickStat(ctx, linkID, browser, device, country, source)
	return err
}

func (s *LinkService) GetStats(ctx context.Context, userID int64, linkID int64) (*model.LinkStats, error) {
	link, err := s.linkRepo.FindByID(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска ссылки: %w", err)
	}
	if link == nil {
		return nil, fmt.Errorf("ссылка не найдена")
	}
	if link.UserID == nil || *link.UserID != userID {
		return nil, fmt.Errorf("нет прав на просмотр статистики")
	}

	return s.linkRepo.GetStats(ctx, linkID)
}

func getCountry(ip string) string {
	if ip == "" || ip == "127.0.0.1" || ip == "::1" {
		return "localhost"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=country", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "unknown"
	}

	var result struct {
		Country string `json:"country"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "unknown"
	}

	if result.Country != "" {
		return result.Country
	}
	return "unknown"
}

func parseUserAgent(ua string) (browser, device string) {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	default:
		browser = "Other"
	}

	switch {
	case strings.Contains(ua, "mobile"):
		device = "Mobile"
	case strings.Contains(ua, "tablet"):
		device = "Tablet"
	default:
		device = "Desktop"
	}

	return
}

func parseReferer(referer string) string {
	if referer == "" {
		return "Прямой"
	}
	referer = strings.ToLower(referer)

	switch {
	case strings.Contains(referer, "twitter.com") || strings.Contains(referer, "x.com"):
		return "Twitter"
	case strings.Contains(referer, "t.me"):
		return "Telegram"
	case strings.Contains(referer, "facebook.com"):
		return "Facebook"
	case strings.Contains(referer, "instagram.com"):
		return "Instagram"
	case strings.Contains(referer, "youtube.com"):
		return "YouTube"
	case strings.Contains(referer, "reddit.com"):
		return "Reddit"
	case strings.Contains(referer, "linkedin.com"):
		return "LinkedIn"
	case strings.Contains(referer, "google.com") || strings.Contains(referer, "google.ru"):
		return "Google"
	default:
		return "Other"
	}
}
