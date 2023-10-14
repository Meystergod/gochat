package client

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host         string
	Port         string
	DatabaseName string
	AuthSource   string
	Username     string
	Password     string
}

func NewMongoConfig(authSource, username, password, host, port, db string) *MongoConfig {
	return &MongoConfig{
		Host:         host,
		Port:         port,
		DatabaseName: db,
		AuthSource:   authSource,
		Username:     username,
		Password:     password,
	}
}

func NewMongoDatabase(ctx context.Context, cfg *MongoConfig) (*mongo.Database, error) {
	var url string
	var anonymous bool
	var client *mongo.Client

	logger := zerolog.Ctx(ctx)

	if cfg.Username == "" || cfg.Password == "" {
		anonymous = true
		url = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}

	clientOptions := options.Client().ApplyURI(url)

	if !anonymous {
		clientOptions.SetAuth(options.Credential{
			AuthSource:  cfg.AuthSource,
			Username:    cfg.Username,
			Password:    cfg.Password,
			PasswordSet: true,
		})
	}

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var connectionError error
	var pingError error

	client, connectionError = mongo.Connect(reqCtx, clientOptions)
	if connectionError != nil {
		return nil, errors.New("failed to connect to database")
	}

	pingError = client.Ping(context.Background(), nil)
	if pingError != nil {
		return nil, errors.New("failed to ping to database")
	}

	logger.Info().Msg("successfully connected to the mongo database")

	return client.Database(cfg.DatabaseName), nil
}
