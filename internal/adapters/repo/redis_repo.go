package repo

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/domain"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/ports"
)

type RedisSymbolRepo struct {
	client *redis.Client
}

func NewRedisSymbolRepo(config config.Config) (ports.ISymbolRepo, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisSymbolRepo{
		client: rdb,
	}, nil
}

func (r *RedisSymbolRepo) Put(data domain.SymbolData) error {
	err := r.client.RPush(context.Background(), data.Ticker, data.CurrentPrice, 0).Err()

	return err
}
