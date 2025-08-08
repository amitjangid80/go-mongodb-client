# ✅ go-mongodb-client

A simple utility Go package for connecting to MongoDB.  
This package provides the following features:

1. Connecting to a MongoDB client.
2. Creating a new database once the MongoDB client is connected.
3. Creating collections inside the database.

---

# 📦 Installation

```bash
go get github.com/amitjangid80/go-mongodb-client@latest
```

# 📦 Usage

## 📦 Creating Mongodb Config

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

### Create Unique Index

```go
func CreateIndex() {
	mongodb_client.CreateUniqueIndex("dbName", "collectionName", "emailId")
}
```

### Create Index

```go
func CreateIndex() {
	// For Single Field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"YOUR_FIELD_NAME": 1},
		Options: options.Index().SetUnique(true),
	}

	// For Multiple Fields
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"YOUR_FIELD_NAME": 1, "ANOTHER_FIELD_NAME": 1},
		Options: options.Index().SetUnique(true),
	}

	mongodb_client.CreateIndex("YOUR_DB_NAME", "YOUR_COLLECTION_NAME", indexModel)
}
```

## 📦 Usage of Base Model and Repository Functions

### 📦 Base DML Model
### You can use this Base model in your own domain model which will by default come with below structure

```go
package mongodb_domain

type BaseDmlModel interface {
	SetCreatedBy(by string)
	SetCreatedOn(on string)
	SetModifiedBy(by string)
	SetModifiedOn(on string)
	GetId() string
	SetId(id string)
}

type DmlModel struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedOn  string `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	CreatedBy  string `json:"createdBy" bson:"createdBy"`
	ModifiedOn string `json:"modifiedOn,omitempty" bson:"modifiedOn,omitempty"`
	ModifiedBy string `json:"modifiedBy" bson:"modifiedBy"`
}

// Implement Auditable interface
func (d *DmlModel) SetCreatedBy(by string)  { d.CreatedBy = by }
func (d *DmlModel) SetCreatedOn(on string)  { d.CreatedOn = on }
func (d *DmlModel) SetModifiedBy(by string) { d.ModifiedBy = by }
func (d *DmlModel) SetModifiedOn(on string) { d.ModifiedOn = on }
func (d *DmlModel) GetId() string           { return d.Id }
func (d *DmlModel) SetId(id string)         { d.Id = id }
```

### Example User Model or Domain

```go
type User struct {
	FirstName      string   `json:"firstName" bson:"firstName"`
	LastName       string   `json:"lastName" bson:"lastName"`
	MobileNumber   string   `json:"mobileNo" bson:"mobileNo"`
	EmailId        string   `json:"emailId" bson:"emailId"`
	Username       string   `json:"username" bson:"username"`
	Password       string   `json:"password" bson:"password"`
	DmlModel       `bson:",inline"` // Add the Base DML Model Here in your domain or model
}
```

## 📦 Repository Functions

### Create or Insert document in a collection

```go
func RegisterUser(user *domain.User) {
	// This function returns the result or error
	result, err := base_repository.CreateRepository[*domain.User]().Create(
		user, 
		"YOUR_DB_NAME", 
		"YOUR_COLLECTION_NAME", 
		"username", // Username who is creating this document to store in createdBy, modifiedBy
	)
}
```

### Update document in a collection

```go
func UpdateUser(user *domain.User) {
	// This function returns the result or error
	result, err := base_repository.UpdateRepository[*domain.User]().Update(
		user, // Include id in user data to update the document
		"YOUR_DB_NAME",
		"YOUR_COLLECTION_NAME",
		"username", // Username who is creating this document to store in modifiedBy
	)
}
```

### Delete document from a collection

```go
func DeleteUser(id string) {
	// This function returns the result or error
	result, err := base_repository.DeleteRepository[*domain.User]().Delete(
		id,
		"YOUR_DB_NAME",
		"YOUR_COLLECTION_NAME"
	)
}
```

### Get data by using Get By Filter Repository from a collection and db

```go
func GetUser(username string, password string) []*domain.User {
	// create options for mongo db query to get the document or documents
	// You can refer mongodb documents for available options
	filterOptions := options.Find()

	// filter for mongodb query to get the document or documents
	filter := bson.M{
		"username": username,
		"password": password,
	}

	// This function will return the result or error
	results, err := base_repository.GetByFilterRepository[*domain.User]().GetByFilter(
		filter,
		filterOptions,
		"YOUR_DB_NAME",
		"YOUR_COLLECTION_NAME",
	)

	if err != nil {
		results = make([]*domain.User, 0)
		log.Printf("Failed to get client from DB: %v", err)
	}

	return results
}
```

### Get data by using Get All Repository from a collection and db

```go
func GetUser() []*domain.User {
	// This function will return the result or error
	results, err := base_repository.GetAllRepository[*domain.User]().GetAll(
		"YOUR_DB_NAME",
		"YOUR_COLLECTION_NAME",
	)

	if err != nil {
		results = make([]*domain.User, 0)
		log.Printf("Failed to get client from DB: %v", err)
	}

	return results
}
```

### Get data by using Get By Id Repository from a collection and db

```go
func GetUserById(id string) []*domain.User {
	// This function will return the result or error
	results, err := base_repository.GetByIdRepository[*domain.User]().GetById(
		id,
		"YOUR_DB_NAME",
		"YOUR_COLLECTION_NAME",
	)

	if err != nil {
		results = make([]*domain.User, 0)
		log.Printf("Failed to get client from DB: %v", err)
	}

	return results
}
```
