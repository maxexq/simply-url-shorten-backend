package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	"url-shortener/handlers"
	"url-shortener/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

var rdb *redis.Client

func init() {
	loadEnv()
	models.ConnectDB()
	setupRedis()
}

func setupRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
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

func main() {
	app := fiber.New()
	app.Post("/shorten", handlers.ShortenURL)
	app.Get("/:code", handlers.RedirectURL)

	log.Fatal(app.Listen(":8080"))
}
