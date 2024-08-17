package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisCl *redis.Client
)

func redisConnect(r redisInfo) {
	redisCl = redis.NewClient(&redis.Options{
		Addr:         r.Host + ":" + strconv.Itoa(r.Port),
		Password:     r.Password,
		DB:           0,
		MinIdleConns: 3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := redisCl.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Redis error: ", err)
	}
}

func RedisByteSETEX(key string, expire int, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return redisCl.Set(ctx, key, value, time.Duration(expire)).Err()
}
