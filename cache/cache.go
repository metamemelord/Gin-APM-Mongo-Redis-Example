package cache

import (
	"context"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/configuration"
	"time"

	"github.com/go-redis/redis/v7"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewClient),
)

type Cache interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte, map[string]interface{}) error
}

type cacheClient struct {
	c *redis.Client
}

func NewClient(config *configuration.Configuration) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Cache.Address,
		Password: config.Cache.Password,
		DB:       config.Cache.DB,
	})

	_, err := client.Ping().Result()
	return &cacheClient{c: client}, err
}

func (c *cacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.c.Get(key)
	return result.Bytes()
}

func (c *cacheClient) Set(ctx context.Context, key string, value []byte, options map[string]interface{}) error {
	if timeout, ok := options["timeout"]; ok {
		return c.c.Set(key, value, time.Duration(timeout.(int))).Err()
	}
	return c.c.Set(key, value, 0).Err()
}
