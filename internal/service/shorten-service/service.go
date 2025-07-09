package shortenservice

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	model "url-shortener/internal/model/shorten-model"
	repo "url-shortener/internal/repository/shorten-repo"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type URLService interface {
	Shorten(url string) (*model.URL, error)
	Resolve(code string) (string, error)
}

type urlService struct {
	repo repo.URLRepository
	rdb  *redis.Client
}

func NewURLService(repo repo.URLRepository, rdb *redis.Client) URLService {
	return &urlService{repo, rdb}
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func (s *urlService) Shorten(url string) (*model.URL, error) {
	shortCode := generateShortCode()
	urlModel := &model.URL{
		ID:          uuid.New().String(),
		OriginalURL: url,
		ShortCode:   shortCode,
	}
	if err := s.repo.SaveURL(urlModel); err != nil {
		return nil, err
	}
	return urlModel, nil
}

func (s *urlService) Resolve(code string) (string, error) {
	ctx := context.Background()

	// Check Redis cache
	if cachedURL, err := s.rdb.Get(ctx, code).Result(); err == nil {
		s.rdb.Incr(ctx, fmt.Sprintf("clicks:%s", code))
		return cachedURL, nil
	}

	// Fallback to DB
	originalURL, err := s.repo.GetOriginalURL(code)
	if err != nil {
		return "", err
	}

	// Cache and increase click count
	s.rdb.Set(ctx, code, originalURL, 24*time.Hour)
	s.rdb.Incr(ctx, fmt.Sprintf("clicks:%s", code))

	return originalURL, nil
}
