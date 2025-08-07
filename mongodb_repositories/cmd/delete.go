package mongodb_repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/amitjangid80/go-mongodb-client/mongodb_client"
	"github.com/amitjangid80/go-mongodb-client/mongodb_domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDeleteRepository[T mongodb_domain.BaseDmlModel] interface {
	Delete(id string, dbName string, collectionName string) (*T, error)
}

type DeleteRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func DeleteRepository[T mongodb_domain.BaseDmlModel]() IDeleteRepository[T] {
	return &DeleteRepositoryImpl[T]{}
}

func (r *DeleteRepositoryImpl[T]) Delete(id, dbName, collectionName string) (*T, error) {
	var deletedDoc T

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("[Delete] Invalid id: %s", id)
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	collection := mongodb_client.GetDb(dbName).Collection(collectionName)
	result := collection.FindOneAndDelete(ctx, filter)

	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("❌ [Delete] No document found for ID: %s", id)
			return nil, err
		}

		log.Printf("❌ [Delete] Error deleting document: %v", err)
		return nil, err
	}

	if err := result.Decode(&deletedDoc); err != nil {
		log.Printf("❌ [Delete] Failed to decode deleted document: %v", err)
		return nil, err
	}

	return &deletedDoc, nil
}
