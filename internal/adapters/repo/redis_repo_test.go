package repo

import (
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestPut_Success(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	repo := &RedisSymbolRepo{client: rdb}

	data := domain.SymbolData{
		Ticker:       "BTC",
		CurrentPrice: 42000.5,
	}

	exactArgs := &redis.XAddArgs{
		Stream: data.Ticker,
		Values: map[string]interface{}{
			"price":     data.CurrentPrice,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	mock.ExpectXAdd(exactArgs).SetVal("1-0")

	err := repo.Put(data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPut_Failure(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	repo := &RedisSymbolRepo{client: rdb}

	data := domain.SymbolData{Ticker: "ETH", CurrentPrice: 3000.0}
	exactArgs := &redis.XAddArgs{
		Stream: data.Ticker,
		Values: map[string]interface{}{
			"price":     data.CurrentPrice,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	mock.ExpectXAdd(exactArgs).SetErr(errors.New("xadd failed"))

	err := repo.Put(data)
	assert.Error(t, err)
	assert.EqualError(t, err, "xadd failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_Success(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	repo := &RedisSymbolRepo{client: rdb}

	entries := []redis.XMessage{
		{ID: "1-0", Values: map[string]interface{}{"price": 100.0}},
		{ID: "2-0", Values: map[string]interface{}{"price": 200.5}},
	}
	mock.ExpectXRangeN("ABC", "-", "+", int64(2)).SetVal(entries)

	result, err := repo.Get("ABC", 2)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 100.0, result[0].CurrentPrice)
	assert.Equal(t, "ABC", result[0].Ticker)
	assert.Equal(t, 200.5, result[1].CurrentPrice)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_Failure(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	repo := &RedisSymbolRepo{client: rdb}

	mock.ExpectXRangeN("XYZ", "-", "+", int64(1)).SetErr(errors.New("xrange error"))

	_, err := repo.Get("XYZ", 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "xrange error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
