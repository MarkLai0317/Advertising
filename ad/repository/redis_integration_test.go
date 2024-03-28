//go:build integration
// +build integration

package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/repository"
	advertisementMock "github.com/MarkLai0317/Advertising/mocks/ad"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RedisIntegrationTestSuite struct {
	testRedisClient *redis.Client
	suite.Suite
}

func TestRedisIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &RedisIntegrationTestSuite{})
}

func (its *RedisIntegrationTestSuite) SetupSuite() {
	// Connect to Redis
	its.testRedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := its.testRedisClient.Ping(context.Background()).Result()
	its.Require().Nil(err, "error connecting to Redis")
}

func (its *RedisIntegrationTestSuite) TearDownSuite() {
	its.Require().Nil(its.testRedisClient.Close(), "error closing Redis connection")
}

func (its *RedisIntegrationTestSuite) SetupTest() {
	its.Require().Nil(its.testRedisClient.FlushDB(context.Background()).Err(), "error flushing Redis database")
}

func (its *RedisIntegrationTestSuite) TearDownTest() {
	its.Require().Nil(its.testRedisClient.FlushDB(context.Background()).Err(), "error flushing Redis database")
}

func (its *RedisIntegrationTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestWithCache" {
		adClient := ad.Client{}
		redisKey, _ := json.Marshal(adClient)
		adSlice := []ad.Advertisement{{}}
		adSliceBytes, _ := json.Marshal(adSlice)
		its.testRedisClient.Set(context.Background(), string(redisKey), adSliceBytes, 0)
	}
}

func (its *RedisIntegrationTestSuite) TestWithCache() {
	adClient := ad.Client{}
	mockRepo := advertisementMock.NewRepository(its.T())
	cacheRepo := repository.NewCacheRepo("localhost:6379", mockRepo)
	adSlice, err := cacheRepo.GetAdvertisements(&adClient, time.Now())
	its.Equal(nil, err, "error getting ads")
	its.Equal([]ad.Advertisement{{}}, adSlice)
}

func (its *RedisIntegrationTestSuite) TestWithoutCache() {

	// setting up input paremeters for GetAdvertisements
	adClient := ad.Client{}
	now := time.Now()
	// create cache repo with injection of mockRepo
	mockRepo := advertisementMock.NewRepository(its.T())
	mockRepo.EXPECT().GetAdvertisements(&adClient, now).Return([]ad.Advertisement{{}}, nil).Times(1)
	cacheRepo := repository.NewCacheRepo("localhost:6379", mockRepo)

	// call GetAdvertisements with adClient
	adSlice, err := cacheRepo.GetAdvertisements(&adClient, now)
	its.Equal(nil, err, "error getting ads")
	its.Equal([]ad.Advertisement{{}}, adSlice)
}
