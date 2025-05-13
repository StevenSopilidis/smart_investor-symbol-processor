package repo

import (
	"context"
	"time"

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
	_, err := r.client.XAdd(
		context.Background(),
		&redis.XAddArgs{
			Stream: data.Ticker,
			Values: map[string]interface{}{
				"price":     data.CurrentPrice,
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
	).Result()

	return err
}

func (r *RedisSymbolRepo) Get(ticker string, ammount int) ([]domain.SymbolData, error) {
	entries, err := r.client.XRangeN(context.Background(), ticker, "-", "+", int64(ammount)).Result()

	if err != nil {
		return nil, err
	}

	result := make([]domain.SymbolData, 0)
	for _, entry := range entries {
		result = append(result, domain.SymbolData{
			CurrentPrice: entry.Values["price"].(float64),
			Ticker:       ticker,
		})
	}

	return result, nil
}
