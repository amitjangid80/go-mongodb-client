package mongodb_query

import (
	"context"
	"log"
	"time"

	"github.com/amitjangid80/go-mongodb-client/mongodb_client"
	"github.com/amitjangid80/go-mongodb-client/mongodb_domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IGetByFilterRepository[T mongodb_domain.BaseDmlModel] interface {
	// If you want multiple docs, change the return type to []T  (see below).
	GetByFilter(filter any, findOpts *options.FindOptions, dbName, collectionName string) ([]T, error)
}

type GetByFilterRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func GetByFilterRepository[T mongodb_domain.BaseDmlModel]() IGetByFilterRepository[T] {
	return &GetByFilterRepositoryImpl[T]{}
}

func (r *GetByFilterRepositoryImpl[T]) GetByFilter(filter any, findOpts *options.FindOptions, dbName, collectionName string) ([]T, error) {
	var results []T = make([]T, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongodb_client.GetDb(dbName).Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		log.Printf("❌ [GetByFilter] error while getting data by filter: %v", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var raw bson.M

		if err := cursor.Decode(&raw); err != nil {
			log.Printf("❌ [GetByFilter] error while decoding data %v", err)
			continue
		}

		var model T

		bsonBytes, _ := bson.Marshal(raw)
		_ = bson.Unmarshal(bsonBytes, &model)

		if idVal, ok := raw["_id"]; ok {
			if objectId, ok := idVal.(primitive.ObjectID); ok {
				model.SetId(objectId.Hex())
			}
		}

		results = append(results, model)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("❌ [GetByFilter] error in mongodb cursor: %s", err)
		return results, err
	}

	return results, nil
}
