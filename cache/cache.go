package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
)

type Config struct {
	Host     string
	Password string
	DB       int
}

type Cache interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte, map[string]interface{}) error
}

type cacheClient struct {
	c *redis.Client
}

func NewClient(config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	return client, err
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
