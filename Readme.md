# go-mongodb-client

A simple utility Go package for connecting to MongoDB.  
This package provides the following features:

1. Connecting to a MongoDB client.
2. Creating a new database once the MongoDB client is connected.
3. Creating collections inside the database.

---

## 📦 Installation

```bash
go get github.com/amitjangid80/go-mongodb-client@latest
```

## 📦 Usage

### You can add this in your main.go file or wherever you are setting up your db connection

```go
import "github.com/amitjangid80/go-mongodb-client/mongodb_client"

func main() {
    mongodbConfig := mongodb_client.MongodbConfig{
        Username: "your_mongo_db_username",
        Password: "your_mongodb_password",
        Host:     "localhost",
        Port:     "your_mongodb_port",
    }

    mongodb_client.ConnectDb(&mongodbConfig)
}
```

### Get database by name

```go
func GetDb() {
    // this function returns *mongo.Database
    database := mongodb_client.GetDb("your_database_name")
}
```

### Create collections

```go
// Create Collections
func CreateCollections(config *config.Config) {
	bgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
	defer cancel()

	collections := []string{
		// list of your collections which you want to create under a database
	}

	// List existing collections
	existingCollections, err := mongodb_client.GetDb().ListCollectionNames(ctx, bson.D{})
	log.Printf("✅ List of Collections: %s", existingCollections)

	if err != nil {
		log.Fatalf("❌ Failed to list collections in DB %s: %v", config.ClientDbName, err)
	}

	existingMap := make(map[string]bool)

	for _, name := range existingCollections {
		existingMap[name] = true
	}

	log.Printf("✅ Collections already present: %v", existingMap)

	// Create only missing collections
	for _, collection := range collections {
		if !existingMap[collection] {
			if err := mongodb_client.GetDb("your_database_name").CreateCollection(ctx, collection); err != nil {
				log.Fatalf("❌ Failed to create collection %s: %v", collection, err)
			}

			log.Printf("✅ Created collection: %s", collection)
		}
	}
}

```
