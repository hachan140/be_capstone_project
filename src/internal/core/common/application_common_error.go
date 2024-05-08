package common

import "errors"

var (
	ErrorConnectRedis       = errors.New("cannot connect to redis")
	ErrorRedisValueNotExist = errors.New("redis key not exist")
)
