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

	// define repo
	dbUrl := os.Getenv("DB_URL")

	// dbUrlArrayStr := os.Getenv("DB_URL_ARRAY")
	// dbUrlArray := strings.Split(dbUrlArrayStr, "|")

	dbTimeoutSecond, err := strconv.Atoi(os.Getenv("DB_TIMEOUT_SECOND"))
	if err != nil {
		log.Fatalf("DB_TIMEOUT format error: %s", err)
	}
	dbRetries, err := strconv.Atoi(os.Getenv("DB_RETRIES"))
	if err != nil {
		log.Fatalf("DB_RETRIES format error: %s", err)
	}

	// repoList := make([]ad.Repository, len(dbUrlArray))

	// for i, dbUrl := range dbUrlArray {
	// 	mongoRepo := repository.NewMongo(dbUrl, time.Duration(dbTimeoutSecond)*time.Second, dbRetries)
	// 	repoList[i] = mongoRepo
	// }

	// loadBalancer := repository.NewLoadBalancerOptions(repoList)

	mongoRepo := repository.NewMongo(dbUrl, time.Duration(dbTimeoutSecond)*time.Second, dbRetries)

	// define usecase service and data transferer
	adService := ad.NewService(mongoRepo, mongoRepo)
	dataTransferer := controller.NewAdDataTransferer()

	adController := controller.NewAdvertisementController(adService, dataTransferer)

	adRouter := router.NewChiAdapter()

	adRouter.Post("/api/v1/ad", adController.CreateAdvertisement)
	adRouter.Get("/api/v1/ad", adController.Advertise)
	adRouter.Get("health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Println("health check")
	})

	err = adRouter.ListenAndServe(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("ListenAndServe error: %s", err)
	}

}
