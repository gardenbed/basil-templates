package client

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// HTTP is a redis.Client that implements the graceful.Client and graceful.Client health.Checker interfaces.
type Redis struct {
	redis.UniversalClient
}

// NewRedis creates a new redis client.
func NewRedis(redisAddress string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{
		UniversalClient: client,
	}
}

// String returns a name for the client.
func (c *Redis) String() string {
	return "redis-client"
}

// Connect opens a long-lived connection to the redis backend.
func (c *Redis) Connect() error {
	ctx := context.Background()
	return c.Ping(ctx).Err()
}

// Disconnect closes the long-lived connection to the redis backend.
func (c *Redis) Disconnect(ctx context.Context) error {
	return c.Close()
}

// HealthCheck checks the health of connection to the redis backend.
func (c *Redis) HealthCheck(ctx context.Context) error {
	return c.Ping(ctx).Err()
}
