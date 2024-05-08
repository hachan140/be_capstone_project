package redis_repo

import (
	"be-capstone-project/src/internal/core/common"
	"be-capstone-project/src/internal/core/dtos"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type IRedisRepository interface {
	GetApplicationConfig(ctx context.Context, env, channel string) (*dtos.FullApplicationInfo, error)
	SetApplicationConfig(ctx context.Context, appInfo dtos.FullApplicationInfo, env, channel string) error
	DeleteApplicationConfig(ctx context.Context, env, channel string)
	Get(c context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(c context.Context, key, field string) (string, error)
	Expire(ctx context.Context, key string, expiration time.Duration)
	//HGetInt Return param value isExist, valueInt, err
	HGetInt(ctx context.Context, key, field string) (int, error)
	GetInt(ctx context.Context, key string) (int, error)
	GetTTLSecondInt(ctx context.Context, key string) int
	GetBytes(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, keys ...string) error
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) (IRedisRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisRepository{client: client}, nil
}

func (r RedisRepository) DeleteApplicationConfig(ctx context.Context, env, channel string) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) GetApplicationConfig(ctx context.Context, env, channel string) (*dtos.FullApplicationInfo, error) {
	var app *dtos.FullApplicationInfo
	res, err := r.client.Get(ctx, fmt.Sprintf(common.ApplicationConfig, env, channel)).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r RedisRepository) SetApplicationConfig(ctx context.Context, appInfo dtos.FullApplicationInfo, env, channel string) error {
	res, err := json.Marshal(appInfo)

	if err != nil {
		return err
	}

	res1 := r.client.Set(ctx, fmt.Sprintf(common.ApplicationConfig, env, channel), res, 0)
	if res1 != nil && res1.Err() != nil {
		return res1.Err()
	}
	return nil
}

func (r RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisRepository) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

func (r RedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) HGetInt(ctx context.Context, key, field string) (int, error) {
	out, err := r.client.HGet(ctx, key, field).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, common.ErrorConnectRedis
		}
		return 0, common.ErrorRedisValueNotExist
	}
	return out, nil
}
func (r RedisRepository) GetInt(ctx context.Context, key string) (int, error) {
	return r.client.Get(ctx, key).Int()
}

func (r RedisRepository) GetTTLSecondInt(ctx context.Context, key string) int {
	return int(r.client.TTL(ctx, key).Val().Seconds())
}

func (r RedisRepository) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r RedisRepository) Get(c context.Context, key string) (string, error) {
	return r.client.Get(c, key).Result()
}

func (r RedisRepository) HGet(c context.Context, key, field string) (string, error) {
	return r.client.HGet(c, key, field).Result()
}

func (r RedisRepository) GetBytes(c context.Context, key string) ([]byte, error) {
	return r.client.Get(c, key).Bytes()
}
