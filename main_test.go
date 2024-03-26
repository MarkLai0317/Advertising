//go:build integration
// +build integration

package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising/ad/repository"
	"github.com/MarkLai0317/Advertising/test_data/integration_test_data/main_test_cases"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MainIntegrationTestSuite struct {
	testMongoClient *mongo.Client
	suite.Suite
}

func TestMainIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &MainIntegrationTestSuite{})
}

const mongoUri = "mongodb://mark:markpwd@localhost:27017"

// set up database for test
func (its *MainIntegrationTestSuite) SetupSuite() {
	mongoClientOptions := options.Client().ApplyURI(mongoUri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mongoClient *mongo.Client
	var err error
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Attempt to connect to MongoDB
		mongoClient, err = mongo.Connect(ctx, mongoClientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB (attempt %d/%d): %s\n", attempt, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Ping the MongoDB server to verify that the client is connected
		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB (attempt %d/%d): %s\n", attempt, maxRetries, err)
			if err := mongoClient.Disconnect(ctx); err != nil {
				log.Printf("Failed to disconnect client: %s\n", err)
			}
			time.Sleep(2 * time.Second)
			continue
		}
		// Connection successful, break out of the retry loop
		break
	}

	if err != nil {
		log.Fatalf("Max retries reached, unable to establish connection to MongoDB: %s\n", err)
	}

	its.testMongoClient = mongoClient

}

// disconnect mongo after test suite complete
func (its *MainIntegrationTestSuite) TearDownSuite() {

	err := its.testMongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("Max retries reached, unable to establish connection to MongoDB: %s\n", err)
	}
}

// clear database before and after every test
func (its *MainIntegrationTestSuite) SetupTest() {

	ctx := context.TODO()
	collection := its.testMongoClient.Database("advertising").Collection("advertisement")

	filter := bson.M{}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
}

func (its *MainIntegrationTestSuite) TearDownTest() {

	ctx := context.TODO()
	collection := its.testMongoClient.Database("advertising").Collection("advertisement")

	filter := bson.M{}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
}

type AdRequest struct {
	Title      string `json:"title"`
	StartAt    string `json:"startAt"`
	EndAt      string `json:"endAt"`
	Conditions struct {
		AgeStart  int      `json:"ageStart"`
		AgeEnd    int      `json:"ageEnd"`
		Genders   []string `json:"gender"`
		Countries []string `json:"country"`
		Platforms []string `json:"platform"`
	} `json:"conditions"`
}

// test create advertisement api
func (its *MainIntegrationTestSuite) TestCreateAdvertisement() {

	//create input body for post request
	inputAdvertisements := AdRequest{
		Title:   "AD 1",
		StartAt: "2023-12-10T03:00:00.000Z",
		EndAt:   "2025-12-31T16:00:00.000Z",
	}
	inputAdvertisements.Conditions.AgeStart = 1
	inputAdvertisements.Conditions.AgeEnd = 100
	inputAdvertisements.Conditions.Genders = []string{"M"}
	inputAdvertisements.Conditions.Countries = []string{"TW", "JP"}
	inputAdvertisements.Conditions.Platforms = []string{"ios", "android"}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(inputAdvertisements)
	if err != nil {
		panic(err)
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", "http://localhost:80/api/v1/ad", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error %s: check test cases", err.Error())
	}

	// statusCode should be 200
	its.Equal(http.StatusOK, resp.StatusCode)

	// connect to DB to see if insert success
	ctx := context.TODO()
	collection := its.testMongoClient.Database("advertising").Collection("advertisement")
	filter := bson.D{{}}
	opts := options.Find()
	// run query on database
	cursor, err := collection.Find(ctx, filter, opts)
	its.Equal(nil, err, "error finding documents, check mongo_test")
	defer cursor.Close(ctx)
	// Iterate through the cursor
	var results []repository.AdvertisementMongo
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("error processing cursor: %s", err.Error())
	}
	its.Equal(nil, err, "error processing cursor, check mongo_test")

	// if the documents insert are successfully inserted by checking the value of each field
	its.Equal(inputAdvertisements.Title, results[0].Title)
	its.Equal(inputAdvertisements.StartAt, results[0].StartAt.Format("2006-01-02T15:04:05.000Z"))
	its.Equal(inputAdvertisements.EndAt, results[0].EndAt.Format("2006-01-02T15:04:05.000Z"))
	its.Equal(inputAdvertisements.Conditions.AgeStart, results[0].Conditions.AgeStart)
	its.Equal(inputAdvertisements.Conditions.AgeEnd, results[0].Conditions.AgeEnd)
	its.Equal(inputAdvertisements.Conditions.Genders, convertEnumSliceToStringSlice(results[0].Conditions.Genders))
	its.Equal(inputAdvertisements.Conditions.Countries, convertEnumSliceToStringSlice(results[0].Conditions.Countries))
	its.Equal(inputAdvertisements.Conditions.Platforms, convertEnumSliceToStringSlice(results[0].Conditions.Platforms))

}

func convertEnumSliceToStringSlice[T ~string](enumSlice []T) []string {
	stringSlice := make([]string, len(enumSlice))
	for i, enumValue := range enumSlice {
		stringSlice[i] = string(enumValue)
	}
	return stringSlice
}

// test get advertisement api
func (its *MainIntegrationTestSuite) TestGetAdvertisement() {

	testCases := main_test_cases.GetAdvertisementsTestCases()

	for name, tc := range testCases {
		its.Run(name, func() {

			its.SetupTest()

			// prepare testData in DB
			collection := its.testMongoClient.Database("advertising").Collection("advertisement")
			_, err := collection.InsertMany(context.TODO(), tc.TestData)
			if err != nil {
				log.Printf("insert document error %s", err.Error())
			}
			its.Equal(nil, err, "prepare data error. pleas check test: TestGetAadvertisements ")

			url := tc.InputUrl
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("error %s: check test cases", err.Error())
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			its.Equal(nil, err)

			var objBody, objExpect interface{}

			// Unmarshal the JSON bytes and string into a Go data structure
			err = json.Unmarshal(bodyBytes, &objBody)
			its.Equal(nil, err)
			err = json.Unmarshal([]byte(tc.Expects.ReturnData), &objExpect)
			its.Equal(nil, err)

			its.TearDownTest()

		})
	}
}
