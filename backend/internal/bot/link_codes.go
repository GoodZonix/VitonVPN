package bot

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LinkCodes struct {
	rdb *redis.Client
}

func NewLinkCodes(rdb *redis.Client) *LinkCodes {
	return &LinkCodes{rdb: rdb}
}

func (s *LinkCodes) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	code := uuid.NewString()
	key := "tg_link:" + code
	if err := s.rdb.Set(ctx, key, userID.String(), ttl).Err(); err != nil {
		return "", err
	}
	return code, nil
}

func (s *LinkCodes) Consume(ctx context.Context, code string) (uuid.UUID, bool, error) {
	key := "tg_link:" + code
	val, err := s.rdb.GetDel(ctx, key).Result()
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

