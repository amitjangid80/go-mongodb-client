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


import (
    "github.com/amitjangid80/go-mongodb-client/mongodb_client"
)

func main() {
    mongodbConfig := mongodb_client.MongodbConfig{
        Username: "your_mongo_db_username",
        Password: "your_mongodb_password",
        Host:     "localhost",
        Port:     "your_mongodb_port",
        DbName:   "your_db_name",
    }

    mongodb_client.ConnectDb(&mongodbConfig)
}
