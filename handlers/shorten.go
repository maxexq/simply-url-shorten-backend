package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	db  *sql.DB
	rdb *redis.Client
)

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func ShortenURL(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "URL is required"})
	}

	shortCode := generateShortCode()
	_, err := db.Exec("INSERT INTO urls (id, original_url, short_code) VALUES (?, ?, ?)", uuid.New().String(), url, shortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to shorten URL"})
	}

	return c.JSON(fiber.Map{"short_url": fmt.Sprintf("http://localhost:8080/%s", shortCode)})
}

func RedirectURL(c *fiber.Ctx) error {
	shortCode := c.Params("code")
	ctx := context.Background()

	if cachedURL, err := rdb.Get(ctx, shortCode).Result(); err == nil {
		log.Println("Cache hit!")
		rdb.Incr(ctx, fmt.Sprintf("clicks:%s", shortCode))
		return c.Redirect(cachedURL, fiber.StatusMovedPermanently)
	}

	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode).Scan(&originalURL)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
	}

	rdb.Set(ctx, shortCode, originalURL, 24*time.Hour)
	rdb.Incr(ctx, fmt.Sprintf("clicks:%s", shortCode))

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}
