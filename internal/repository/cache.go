package repository

import (
	"context"
	"os"
	"time"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"github.com/redis/go-redis/v9"
)

const cacheCaller = packageCaller + ".Cache"

type Cache struct {
	client *redis.Client
}

func NewCacheRepo(logger *utils.Logger) ports.CacheRepo {
	host := os.Getenv("DRAGONFLYDB_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_HOST is not set")
	}
	port := os.Getenv("DRAGONFLYDB_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_PORT is not set")
	}
	pass := os.Getenv("DRAGONFLYDB_PASSWORD")
	if pass == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_PASSWORD is not set")
	}
	db := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})
	_, err := db.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to connect to Dragonfly DB: %v", err)
	}
	return &Cache{
		client: db,
	}
}

func (c *Cache) AddJwtToBlacklist(ctx context.Context, token string, expire time.Duration) (err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".AddJwtToBlacklist", err) }()
	if err := c.client.Set(ctx, token, "true", expire).Err(); err != nil {
		return err
	}
	return
}

func (c *Cache) IsJwtInBlacklist(ctx context.Context, token string) (isIn bool, err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".IsJwtInBlacklist", err) }()
	if err := c.client.Get(ctx, token).Scan(&isIn); err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Cache) SetOtpToken(ctx context.Context, identity, token string, otpType ports.OtpType, expire time.Duration, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".SetOtpToken", err) }()
	if err := c.client.Set(ctx, identity+string(userType)+string(otpType), token, expire).Err(); err != nil {
		return err
	}
	return
}

func (c *Cache) GetOtpToken(ctx context.Context, identity string, otpType ports.OtpType, userType ports.UserType) (token string, err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".VerifyOtpToken", err) }()
	if err := c.client.Get(ctx, identity+string(userType)+string(otpType)).Scan(&token); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return

}

func (c *Cache) DeleteOtpToken(ctx context.Context, identity string, otpType ports.OtpType, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".DeleteOtpToken", err) }()
	if err := c.client.Del(ctx, identity+string(userType)+string(otpType)).Err(); err != nil {
		return err
	}
	return
}

func (c *Cache) SetOtpKey(ctx context.Context, identity, key string, expire time.Duration, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".SetOtpKey", err) }()
	if err := c.client.Set(ctx, identity+string(userType)+"key", key, expire).Err(); err != nil {
		return err
	}
	return
}

func (c *Cache) GetOtpKey(ctx context.Context, identity string, userType ports.UserType) (key string, err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".GetOtpKey", err) }()
	if err := c.client.Get(ctx, identity+string(userType)+"key").Scan(&key); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return
}

func (c *Cache) DeleteOtpKey(ctx context.Context, identity string, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(cacheCaller+".DeleteOtpKey", err) }()
	if err := c.client.Del(ctx, identity+string(userType)+"key").Err(); err != nil {
		return err
	}
	return
}

func (c *Cache) Close() error {
	return c.client.Close()
}
