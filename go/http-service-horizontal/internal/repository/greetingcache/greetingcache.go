package greetingcache

import (
	"context"
	"errors"
	"time"

	"github.com/gardenbed/basil/graceful"
	"github.com/gardenbed/basil/health"
	"github.com/go-redis/redis/v8"
)

// Repository is the interface for interacting with the data store.
type Repository interface {
	graceful.Client
	health.Checker
	Store(ctx context.Context, lang, greeting string) error
	Lookup(ctx context.Context, lang string) (string, error)
}

type redisClient interface {
	Close() error
	Ping(context.Context) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
}

// repository implements the Repository interface.
type repository struct {
	client redisClient
}

// NewRepository creates a new repository.
func NewRepository(redisAddress string) (Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &repository{
		client: client,
	}, nil
}

// String returns a name for the repository.
func (r *repository) String() string {
	return "greetingcache-repository"
}

// Connect opens a long-lived connection to the repository backend.
func (r *repository) Connect() error {
	ctx := context.Background()
	return r.client.Ping(ctx).Err()
}

// Disconnect closes the long-lived connection to the repository backend.
func (r *repository) Disconnect(ctx context.Context) error {
	return r.client.Close()
}

// HealthCheck checks the health of connection to the repository backend.
func (r *repository) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Store
func (r *repository) Store(ctx context.Context, lang, greeting string) error {
	if lang == "" {
		return errors.New("no language code")
	}

	if greeting == "" {
		return errors.New("no greeting value")
	}

	return r.client.Set(ctx, lang, greeting, 0).Err()
}

// Lookup
func (r *repository) Lookup(ctx context.Context, lang string) (string, error) {
	if lang == "" {
		return "", errors.New("no language code")
	}

	return r.client.Get(ctx, lang).Result()
}
