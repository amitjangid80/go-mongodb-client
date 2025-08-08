package mongodb_repositories

import (
	"context"
	"log"
	"time"

	"github.com/amitjangid80/go-mongodb-client/mongodb_client"
	"github.com/amitjangid80/go-mongodb-client/mongodb_domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUpdateRepository[T mongodb_domain.BaseDmlModel] interface {
	Update(data T, dbName string, collectionName string, createdBy string) (*T, error)
}

type UpdateRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func UpdateRepository[T mongodb_domain.BaseDmlModel]() IUpdateRepository[T] {
	return &UpdateRepositoryImpl[T]{}
}

// Update implements IUpdateRepository.
func (g *UpdateRepositoryImpl[T]) Update(data T, dbName string, collectionName string, createdBy string) (*T, error) {
	var updatedDoc T

	bgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
	defer cancel()

	data.SetModifiedBy(createdBy)
	data.SetModifiedOn(time.Now().UTC().Format(time.RFC3339))

	objectId, err := primitive.ObjectIDFromHex(data.GetId())

	if err != nil {
		log.Printf("[Update] Invalid object id: %s", data.GetId())
		return nil, err
	}

	// Convert the data struct to a BSON map
	marshalledData, err := bson.Marshal(data)

	if err != nil {
		log.Printf("[Update] Failed to marshall data: %v", data)
		return nil, err
	}

	// Unmarshal to map and remove empty values
	var updateBsonMap bson.M

	if err := bson.Unmarshal(marshalledData, &updateBsonMap); err != nil {
		log.Printf("[Update] Failed to unmarshal to map: %v", err)
		return nil, err
	}

	// Remove fields you don't want to update if they are empty (optional)
	for key, value := range updateBsonMap {
		if value == nil || value == "" {
			delete(updateBsonMap, key)
		}
	}

	// ✅ prevent trying to update _id
	delete(updateBsonMap, "_id")

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": updateBsonMap}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	collection := mongodb_client.GetDb(dbName).Collection(collectionName)
	updateResult := collection.FindOneAndUpdate(ctx, filter, update, options)

	if err := updateResult.Decode(&updatedDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("❌ [Update] No document found with id: %s", data.GetId())
			return nil, err
		}

		log.Printf("❌ [Update] Error decoding updated document: %v", err)
		return nil, err
	}

	return &updatedDoc, nil
}
