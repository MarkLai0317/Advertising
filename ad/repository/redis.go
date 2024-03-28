package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	redisClient *redis.Client
	mainRepo    ad.Repository
}

func NewCacheRepo(host string, mainRepo ad.Repository) *CacheRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &CacheRepo{redisClient: rdb, mainRepo: mainRepo}
}

func (c *CacheRepo) CreateAdvertisement(ad *ad.Advertisement) error {

	return c.mainRepo.CreateAdvertisement(ad)
}

func (c *CacheRepo) GetAdvertisements(client *ad.Client, now time.Time) ([]ad.Advertisement, error) {

	redisKey, _ := json.Marshal(client)

	cachedAdSlice, err := c.redisClient.Get(context.Background(), string(redisKey)).Result()
	if err == redis.Nil {
		adSlice, err := c.mainRepo.GetAdvertisements(client, now)
		if err != nil {
			return nil, fmt.Errorf("error getting ads from main repo: %w", err)
		}
		adSliceBytes, _ := json.Marshal(adSlice)
		if err := c.redisClient.Set(context.Background(), string(redisKey), adSliceBytes, 0).Err(); err != nil {
			return nil, fmt.Errorf("error setting ads in redis: %w", err)
		}
		return adSlice, nil
	} else if err != nil {
		return nil, fmt.Errorf("error checking if ads exist in redis: %w", err)
	}

	var adSlice []ad.Advertisement
	if err := json.Unmarshal([]byte(cachedAdSlice), &adSlice); err != nil {
		return nil, fmt.Errorf("error unmarshalling cached ads: %w", err)
	}
	return adSlice, nil
}
