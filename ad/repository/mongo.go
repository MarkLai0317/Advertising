package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MarkLai0317/Advertising/ad"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	mongoClient     *mongo.Client
	dbName          string
	readCollection  string
	writeCollection string
}

func NewMongo(uri string, dbName string, writeCollection string, readCollection string, connectTimeout time.Duration, maxRetries int) *Mongo {
	// set mongo connection options and timeout
	mongoClientOptions := options.Client().ApplyURI(uri)
	mongoClientOptions.SetMaxPoolSize(0)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	var mongoClient *mongo.Client
	var err error

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

	return &Mongo{mongoClient: mongoClient, dbName: dbName, writeCollection: writeCollection, readCollection: readCollection}

}

// Advertisement entity in mongo with bson tag
// which can decouple domain entity from external service such as mongodb bson format
// make domain logic not depend on mongo db format
type AdvertisementMongo struct {
	Title      string          `bson:"title"`
	StartAt    time.Time       `bson:"startAt"`
	EndAt      time.Time       `bson:"endAt"`
	Conditions ConditionsMongo `bson:"conditions"`
}

type ConditionsMongo struct {
	AgeStart  int               `bson:"ageStart"`
	AgeEnd    int               `bson:"ageEnd"`
	Genders   []ad.GenderType   `bson:"genders"`
	Countries []ad.CountryCode  `bson:"countries"`
	Platforms []ad.PlatformType `bson:"platforms"`
}

func (m *Mongo) CreateAdvertisement(advertisement *ad.Advertisement) error {
	advertisementMongo := AdvertisementMongo{
		Title:   advertisement.Title,
		StartAt: advertisement.StartAt,
		EndAt:   advertisement.EndAt,
		Conditions: ConditionsMongo{
			AgeStart:  advertisement.Conditions.AgeStart,
			AgeEnd:    advertisement.Conditions.AgeEnd,
			Genders:   advertisement.Conditions.Genders,
			Countries: advertisement.Conditions.Countries,
			Platforms: advertisement.Conditions.Platforms,
		},
	}

	collection := m.mongoClient.Database(m.dbName).Collection(m.writeCollection)
	result, err := collection.InsertOne(context.TODO(), advertisementMongo)
	if err != nil {
		return fmt.Errorf("error inserting advertisement: %w", err)
	}
	log.Printf("Inserted document with _id: %v\n", result.InsertedID)

	return nil

}

func (m *Mongo) GetAdvertisements(client *ad.Client, now time.Time) ([]ad.Advertisement, error) {

	collection := m.mongoClient.Database(m.dbName).Collection(m.readCollection)

	// Define your query using bson.D to ensure order
	ctx := context.TODO()
	filter := buildMongoQuery(client, now)

	opts := options.Find().
		SetSort(bson.D{{"endAt", 1}}).
		SetSkip(int64(client.Offset)).
		SetLimit(int64(client.Limit)).
		SetProjection(bson.D{{"title", 1}, {"endAt", 1}, {"_id", 0}})

	// Execute the query
	cursor, err := collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, fmt.Errorf("error when find advertisement: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor
	results := []AdvertisementMongo{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error converting query result to AdvertisementMongo slice: %w", err)
	}

	// convert AdvertisementMongo to domain object slice []Advertisement
	adSlice := make([]ad.Advertisement, len(results))
	for i, result := range results {
		adSlice[i] = ad.Advertisement{
			Title: result.Title,
			EndAt: result.EndAt,
		}
	}
	return adSlice, nil
}

// create query if the parama exist
func buildMongoQuery(client *ad.Client, now time.Time) bson.D {

	query := bson.D(make([]bson.E, 0, 5))
	if !client.CountryMissing {
		countryQuery := bson.E{"conditions.countries", string(client.Country)}
		query = append(query, countryQuery)
	}

	if !client.AgeMissing {
		ageStartQuery := bson.E{"conditions.ageStart", bson.D{{"$lte", client.Age}}}
		ageEndQuery := bson.E{"conditions.ageEnd", bson.D{{"$gte", client.Age}}}
		query = append(query, ageStartQuery)
		query = append(query, ageEndQuery)
	}

	if !client.GenderMissing {
		genderQuery := bson.E{"conditions.genders", string(client.Gender)}
		query = append(query, genderQuery)
	}

	if !client.PlatformMissing {
		platformQuery := bson.E{"conditions.platforms", string(client.Platform)}
		query = append(query, platformQuery)
	}

	return query
}

// func constructClientQuery(client *ad.Client){

// }

// func ageQuery(missing bool) bson.D {
// 	if missing {
// 		return bson.D{{"$exists", true}}
// 	}
// 	return bson.D{{"$lte", now}}
// }

// func AgeQuery(age int) bson.D{

// 	if

// }
