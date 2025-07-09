package web

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitEnv() {
	var envPath = os.Getenv("DOTENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}
	err := godotenv.Load(envPath)
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
}
