package session_store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
	"github.com/go-redis/redis/v8"
)

type RedisSessionStore struct {
	client redis.Client
}

func NewSessionStore(conf config.RedisConfig) (*RedisSessionStore, error) {
	clientOpts := redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.Db,
	}

	ctx := context.Background()
	client := redis.NewClient(&clientOpts)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %s", err.Error())
	}

	return &RedisSessionStore{client: *client}, nil
}

func (st *RedisSessionStore) Get(key int64) *session.Session {
	ctx := context.Background()

	result, err := st.client.Get(ctx, strconv.FormatInt(key, 10)).Result()
	if err != nil {
		return nil
	}

	var storedSession session.Session
	if err := json.Unmarshal([]byte(result), &storedSession); err != nil {
		return nil
	}

	return &storedSession
}

func (st *RedisSessionStore) Set(key int64, ssn *session.Session) {
	ctx := context.Background()

	raw, err := json.Marshal(&ssn)
	if err != nil {
		panic(err)
	}

	st.client.Set(ctx, strconv.FormatInt(key, 10), raw, 0)
}

func (st *RedisSessionStore) Delete(key int64) {
	ctx := context.Background()

	st.client.Del(ctx, strconv.FormatInt(key, 10))
}
