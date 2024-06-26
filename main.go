package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/controller"
	"github.com/MarkLai0317/Advertising/ad/repository"
	"github.com/MarkLai0317/Advertising/internal/router"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// get env for mongoRepo
	dbUrl := os.Getenv("DB_URL")
	dbTimeoutSecond, err := strconv.Atoi(os.Getenv("DB_TIMEOUT_SECOND"))
	if err != nil {
		log.Fatalf("DB_TIMEOUT format error: %s", err)
	}
	dbRetries, err := strconv.Atoi(os.Getenv("DB_RETRIES"))
	if err != nil {
		log.Fatalf("DB_RETRIES format error: %s", err)
	}
	// define repository
	writeCollection := os.Getenv("WRITE_COLLECTION")
	readCollection := os.Getenv("READ_COLLECTION")
	dbName := os.Getenv("DB_NAME")
	mongoRepo := repository.NewMongo(dbUrl, dbName, writeCollection, readCollection, time.Duration(dbTimeoutSecond)*time.Second, dbRetries)
	var usedRepo ad.Repository
	usedRepo = mongoRepo
	if os.Getenv("USE_CACHE") == "TRUE" {
		redisHost := os.Getenv("REDIS_HOST")
		redisPoolSize, err := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
		if err != nil {
			log.Fatalf("REDIS_POOL_SIZE format error: %s", err)
		}
		cacheRepo := repository.NewCacheRepo(redisHost, redisPoolSize, mongoRepo)
		usedRepo = cacheRepo
	}

	// define usecase service and data transferer
	adService := ad.NewService(usedRepo)
	dataTransferer := controller.NewAdDataTransferer()
	// inject to controller
	adController := controller.NewAdvertisementController(adService, dataTransferer)
	// define router
	adRouter := router.NewChiRouter()
	defineAPI(adRouter, adController)

}

// define route for API
func defineAPI(adRouter router.WebFramework, adController *controller.Controller) {
	adRouter.Post("/api/v1/ad", adController.CreateAdvertisement)
	adRouter.Get("/api/v1/ad", adController.Advertise)
	adRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Println("health check")
	})

	err := adRouter.ListenAndServe(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("ListenAndServe error: %s", err)
	}
}
