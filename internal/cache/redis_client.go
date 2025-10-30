package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(addr, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis")
}

// Helper functions
func Set(key string, value string, ttl time.Duration) {
	if err := Rdb.Set(Ctx, key, value, ttl).Err(); err != nil {
		log.Printf("❌ Redis SET error: %v", err)
	}
}

func Get(key string) (string, bool) {
	val, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		log.Printf("❌ Redis GET error: %v", err)
		return "", false
	}
	return val, true
}
