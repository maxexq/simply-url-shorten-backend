package drivers

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnection driver
func MongoConnection(dsn string) (mongo.Client, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	var client *mongo.Client

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		logrus.Fatalf("can't connect mongodb : %s", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Fatalf("can't connect ping mongodb : %s", err)
	}

	logrus.Info("connect mongodb success")

	return *client, ctx
}
