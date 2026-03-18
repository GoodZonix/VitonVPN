package bot

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CabinetTokens struct {
	rdb *redis.Client
}

func NewCabinetTokens(rdb *redis.Client) *CabinetTokens {
	return &CabinetTokens{rdb: rdb}
}

func (s *CabinetTokens) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	token := uuid.NewString()
	key := "cabinet:" + token
	if err := s.rdb.Set(ctx, key, userID.String(), ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (s *CabinetTokens) Get(ctx context.Context, token string) (uuid.UUID, bool, error) {
	key := "cabinet:" + token
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return uuid.Nil, false, nil
	}
	if err != nil {
		return uuid.Nil, false, err
	}
	uid, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, false, nil
	}
	return uid, true, nil
}

