package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

var (
	db  *sql.DB
	rdb *redis.Client
)

func init() {
	loadEnv()
	connectDB()
	setupRedis()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectDB() {
	var err error
	db, err = sql.Open("sqlite", "./urls.db")
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table", err)
	}
}

func setupRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
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

func shortenURL(c *fiber.Ctx) error {
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

func redirectURL(c *fiber.Ctx) error {
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

func main() {
	app := fiber.New()

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/shorten", shortenURL)
	app.Get("/:code", redirectURL)

	log.Fatal(app.Listen(":8080"))

}
