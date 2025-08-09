package mongodb_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// GetClient returns the mongodb client.
func GetClient() *mongo.Client {
	return client
}

// GetDb returns the mongodb client.
func GetDb(dbName string) *mongo.Database {
	return client.Database(dbName)
}

// Connect Mongodb function will initialize and connect to mongodb based on the URL, Port and Host passed via config
func ConnectDb(config *MongodbConfig) {
	log.Printf("✅ Connecting to MongoDB on port: %s...", config.Port)

	// Get MongoDB URI from environment variable if set, otherwise use default
	mongoDbUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Username, config.Password, config.Host, config.Port)

	clientOpts := options.Client().ApplyURI(mongoDbUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOpts)

	if err != nil {
		log.Printf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("❌ Failed to ping MongoDB: %v", err)
	}

	log.Printf("✅ Connected to MongoDB on port: %s", config.Port)
}

// Create Collection
func CreateCollection(dbName string, collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List existing collections
	existingCollections, err := GetDb(dbName).ListCollectionNames(ctx, bson.D{})
	log.Printf("List of Collections in %s: %v", dbName, existingCollections)

	if err != nil {
		log.Printf("❌ Failed to list collections in DB %s: %v", dbName, err)
	}

	existingMap := make(map[string]bool)

	for _, name := range existingCollections {
		existingMap[name] = true
	}

	// Create only missing collections
	if !existingMap[collectionName] {
		if err := GetDb(dbName).CreateCollection(ctx, collectionName); err != nil {
			log.Printf("❌ Failed to create collection %s: %v", collectionName, err)
		}

		log.Printf("✅ Created collection: %s", collectionName)
	} else {
		log.Printf("✅ Collection already present: %s", collectionName)
	}
}

// Create Collections
func CreateCollections(dbName string, collections []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := GetDb(dbName)

	// List existing collections
	existingCollections, err := db.ListCollectionNames(ctx, bson.D{})
	log.Printf("List of Collections: %s", existingCollections)

	if err != nil {
		log.Printf("Failed to list collections in DB %s: %v", dbName, err)
	}

	existingMap := make(map[string]bool)

	for _, name := range existingCollections {
		existingMap[name] = true
	}

	log.Printf("Collections already present: %v", existingMap)

	// Create only missing collections
	for _, collection := range collections {
		if !existingMap[collection] {
			if err := db.CreateCollection(ctx, collection); err != nil {
				log.Printf("Failed to create collection %s: %v", collection, err)
			}

			log.Printf("Created collection: %s", collection)
		}
	}
}

func CreateIndex(dbName string, collectionName string, indexModel mongo.IndexModel) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetDb(dbName).Collection(collectionName).Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		log.Printf("❌ Failed to create index for collection: %s in db: %s", collectionName, dbName)
		log.Printf("❌ Error while creating creating index: %v", err)
	} else {
		log.Println("✅ Index Created Successfully")
	}
}

func CreateUniqueIndex(dbName string, collectionName string, field string) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetDb(dbName).Collection(collectionName).Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		log.Printf("❌ Failed to create index for collection: %s in db: %s", collectionName, dbName)
		log.Printf("❌ Error while creating creating index: %v", err)
	} else {
		log.Println("✅ Index Created Successfully")
	}
}
