package nosql

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
type MongoConfig struct {
	Address     []string
	MaxPoolSize uint64
}

func (cfg MongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s", strings.Join(cfg.Address, ","))
}

func (cfg *MongoConfig) WithMongoDefault() *MongoConfig {
	if cfg == nil {
		return nil
	}

	c := *cfg

	if cfg.Address == nil || len(cfg.Address) == 0 {
		c.Address = []string{"127.0.0.1:27017"}
	}

	// fuck here
	if cfg.MaxPoolSize == 0 {
		c.MaxPoolSize = 50
	}

	return &c
}

func NewMongo(ctx context.Context, cfg *MongoConfig) (client *mongo.Client, err error) {
	c := cfg.WithMongoDefault()

	if client, err = mongo.Connect(ctx,
		options.Client().
			ApplyURI(c.URI()).
			// default: 100, 0: math.MaxInt64
			SetMaxPoolSize(c.MaxPoolSize)); err != nil {
		return
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	return
}
