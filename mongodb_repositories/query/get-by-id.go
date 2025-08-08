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

type IGetByIdRepository[T mongodb_domain.BaseDmlModel] interface {
	GetById(id string, dbName string, collectionName string, createdBy string) (*T, error)
}

type GetByIdRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func GetByIdRepository[T mongodb_domain.BaseDmlModel]() IGetByIdRepository[T] {
	return &GetByIdRepositoryImpl[T]{}
}

// GetById implements GetByIdRepositoryImpl.
func (g *GetByIdRepositoryImpl[T]) GetById(id string, dbName string, collectionName string, createdBy string) (*T, error) {
	bgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("❌ [GetById] Invalid ID format: %s", err)
		return nil, err
	}

	filter := bson.M{"_id": objectId, "createdBy": createdBy}
	collection := mongodb_client.GetClient().Database(dbName).Collection(collectionName)

	var result T
	err = collection.FindOne(ctx, filter).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("❌ [GetById] mongodb Error: Data Not found: %s", err)
		return nil, err
	}

	if err != nil {
		log.Printf("❌ [GetById] Data Not found: %s", err)
		return nil, err
	}

	result.SetId(id)

	return &result, nil
}
