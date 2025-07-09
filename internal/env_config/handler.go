package envconfig

import (
	"url-shortener/pkg/web"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MongoDBUrl string `env:"MONGO_URL,required"`
	RedisHost  string `env:"REDIS_HOST,required"`
}

func ParseEnv() Config {
	web.InitEnv()
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("Parse Env : %v", err.Error())
		panic(err)
	}
	return cfg
}
