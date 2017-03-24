package util

import (
	"errors"
	"time"

	"apiserver/pkg/configz"

	"gopkg.in/redis.v5"
)

var (
	client *redis.Client
)

func init() {
	addr := configz.GetString("redis", "address", "0.0.0.0:6379")
	password := configz.GetString("redis", "password", "")
	db := configz.MustInt("redis", "db", 0)
	poolSize := configz.MustInt("redis", "poolSize", 10)
	if configz.MustBool("redis", "requiered_password", false) {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			PoolSize: poolSize,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			DB:       db,
			PoolSize: poolSize,
		})
	}
}

func Set(key string, val interface{}, expiredTime int) (string, error) {
	result, err := client.Set(key, val, time.Minute*30).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func Get(key string) (string, error) {
	result, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("key's value is not exsit")
	}
	return result, nil
}

func MutilSet() {
	// this.Client.Pipelined(fn).
}
