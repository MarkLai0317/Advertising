//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/repository"
	"github.com/MarkLai0317/Advertising/test_data/integration_test_data/ad/repository/mongo_test_cases"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoIntegrationTestSuite struct {
	testMongoClient *mongo.Client
	suite.Suite
}

func TestMongoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &MongoIntegrationTestSuite{})
}

const mongoUri = "mongodb://mark:markpwd@localhost:27017"
const writeCollection = "all_advertisement"
const readCollection = "active_advertisement"
const dbName = "advertising"

func (its *MongoIntegrationTestSuite) SetupSuite() {
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
func (its *MongoIntegrationTestSuite) TearDownSuite() {

	err := its.testMongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("Max retries reached, unable to establish connection to MongoDB: %s\n", err)
	}
}

// clear database before and after every test
func (its *MongoIntegrationTestSuite) SetupTest() {

	ctx := context.TODO()
	collection := its.testMongoClient.Database(dbName).Collection(writeCollection)

	filter := bson.M{}

	_, err := collection.DeleteMany(ctx, filter)

	collection = its.testMongoClient.Database(dbName).Collection(readCollection)

	filter = bson.M{}

	_, err = collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
}

func (its *MongoIntegrationTestSuite) TearDownTest() {

	ctx := context.TODO()
	collection := its.testMongoClient.Database(dbName).Collection(writeCollection)

	filter := bson.M{}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	collection = its.testMongoClient.Database(dbName).Collection(readCollection)
	filter = bson.M{}
	_, err = collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
}

func (its *MongoIntegrationTestSuite) TestCreateAdvertisement() {

	//  init mongoRepo object

	mongoRepo := repository.NewMongo(mongoUri, dbName, writeCollection, readCollection, 10*time.Second, 3)
	inputAdvertisements := []ad.Advertisement{
		{
			Title:   "integration test",
			StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndAt:   time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC),
			Conditions: ad.Conditions{
				AgeStart:  1,
				AgeEnd:    100,
				Genders:   []ad.GenderType{"M"},
				Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
				Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
			},
		},
	}

	// test CreateAdvertisement
	for _, advertisement := range inputAdvertisements {
		err := mongoRepo.CreateAdvertisement(&advertisement)
		its.Equal(nil, err, "error creating advertisement")
		if err != nil {
			log.Printf("err creating advertisement: %s", err.Error())
		}
	}

	// check if insert successfully by query the database to see the documents are the same as previous created
	ctx := context.TODO()
	collection := its.testMongoClient.Database(dbName).Collection(writeCollection)
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

	// if the documents insert are successfully inserted by checking len of slice and the value of each field
	its.Equal(len(inputAdvertisements), len(results), "ad insert should be the same as ad find()")
	if len(inputAdvertisements) == len(results) {

		for i, result := range results {
			its.Equal(inputAdvertisements[i].Title, result.Title, fmt.Sprintf("Title should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].StartAt, result.StartAt, fmt.Sprintf("StartAt should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].EndAt, result.EndAt, fmt.Sprintf("EndAt should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].Conditions.AgeStart, result.Conditions.AgeStart, fmt.Sprintf("Condition.AgeStart should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].Conditions.AgeEnd, result.Conditions.AgeEnd, fmt.Sprintf("Condition.AgeEnd should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].Conditions.Genders, result.Conditions.Genders, fmt.Sprintf("Condition.Genders should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].Conditions.Countries, result.Conditions.Countries, fmt.Sprintf("Condition.Countries should be the same in []Advertisement %d", i))
			its.Equal(inputAdvertisements[i].Conditions.Platforms, result.Conditions.Platforms, fmt.Sprintf("Condition.Platforms should be the same in []Advertisement %d", i))
		}
	}

	// mongoRepo.CreateAdvertisement
}

func (its *MongoIntegrationTestSuite) TestGetAdvertisements() {
	testCases := mongo_test_cases.GetAdvertisementsTestCases()

	for name, tc := range testCases {
		its.Run(name, func() {

			its.SetupTest()
			// prepare testData in DB
			collection := its.testMongoClient.Database(dbName).Collection(readCollection)
			_, err := collection.InsertMany(context.TODO(), tc.TestData)
			if err != nil {
				log.Printf("insert document error %s", err.Error())
			}
			its.Equal(nil, err, "prepare data error. pleas check test: TestGetAadvertisements ")

			// run tested function
			mongoRepo := repository.NewMongo(mongoUri, dbName, writeCollection, readCollection, 10*time.Second, 3)
			advertisementSlice, err := mongoRepo.GetAdvertisements(&tc.Input.Client, tc.Input.Now)

			its.Equal(tc.Expects.ExpectError, err, "GetAdvertisements error not the same")
			its.Equal(tc.Expects.ReturnData, advertisementSlice)

			its.TearDownTest()

		})
	}
}
