package mongodb_repositories

import (
	"context"
	"log"
	"time"

	"github.com/amitjangid80/go-mongodb-client/mongodb_client"
	"github.com/amitjangid80/go-mongodb-client/mongodb_domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICreateRepository[T mongodb_domain.BaseDmlModel] interface {
	Create(data T, dbName string, collectionName string, username string) (*T, error)
}

type CreateRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func CreateRepository[T mongodb_domain.BaseDmlModel]() ICreateRepository[T] {
	return &CreateRepositoryImpl[T]{}
}

// Create implements ICreateRepository.
func (g *CreateRepositoryImpl[T]) Create(data T, dbName string, collectionName string, username string) (*T, error) {
	bgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
	defer cancel()

	now := time.Now().UTC().Format(time.RFC3339)

	data.SetCreatedBy(username)
	data.SetCreatedOn(now)
	data.SetModifiedBy(username)
	data.SetModifiedOn(now)

	collection := mongodb_client.GetDb(dbName).Collection(collectionName)
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		log.Printf("❌ [Create] Failed to insert data in collection '%s', error: %v", collectionName, err)
		return nil, err
	}

	// Convert ObjectID to hex string and set in response
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		data.SetId(oid.Hex())
	}

	return &data, nil
}
