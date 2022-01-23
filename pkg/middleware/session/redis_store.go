package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisSessionStore struct {
	client redis.Client
}

func NewRedisSessionStore(opts redis.Options) (*RedisSessionStore, error) {
	ctx := context.Background()
	client := redis.NewClient(&opts)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %s", err.Error())
	}

	return &RedisSessionStore{client: *client}, nil
}

func (st *RedisSessionStore) Get(key int64) *Session {
	ctx := context.Background()

	result, err := st.client.Get(ctx, strconv.FormatInt(key, 10)).Result()
	if err != nil {
		return nil
	}

	var storedSession Session
	if err := json.Unmarshal([]byte(result), &storedSession); err != nil {
		return nil
	}

	return &storedSession
}

func (st *RedisSessionStore) Set(key int64, ssn *Session) {
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

func (st *RedisSessionStore) Stop() error {
	return st.client.Close()
}
