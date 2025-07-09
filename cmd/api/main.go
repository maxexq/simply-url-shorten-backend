package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/pkg/drivers"
	"url-shortener/pkg/routing"

	envconfig "url-shortener/internal/env_config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"

	handler "url-shortener/internal/handler/shorten-handler"
	repository "url-shortener/internal/repository/shorten-repo"
	service "url-shortener/internal/service/shorten-service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	logrus.Infof("TZ=%v", time.Local.String())

	config := envconfig.ParseEnv()

	app := fiber.New()

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	var (
		// Drivers
		dbmg, dbmgCtx = drivers.MongoConnection(config.MongoDBUrl)
		rdb           = drivers.NewRedis(
			config.RedisHost,
			config.RedisHost,
		)
		dbConn = drivers.ConnectDB()

		newFiber = routing.InitFiber()
		f, _     = newFiber.InitFiberMiddleware()
	)

	repo := repository.NewURLRepository(dbConn)
	svc := service.NewURLService(repo, rdb)
	h := handler.NewURLHandler(svc)

	f.Post("/shorten", h.ShortenURL)
	f.Get("/:code", h.RedirectURL)

	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logrus.Info("Gracefully shutting down...")
		f.Shutdown()
	}()

	if err := f.Listen(":" + os.Getenv("PORT")); err != nil {
		logrus.Fatalf("shutting down the server : %s", err)
	}

	defer func() {
		if err := dbmg.Disconnect(dbmgCtx); err != nil {
			panic(err)
		}
	}()
}
