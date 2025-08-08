package mongodb_repositories

import (
	"context"
	"log"
	"time"

	"github.com/amitjangid80/go-mongodb-client/mongodb_client"
	"github.com/amitjangid80/go-mongodb-client/mongodb_domain"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type IGetAllRepository[T mongodb_domain.BaseDmlModel] interface {
	GetAll(dbName string, collectionName string, createdBy string) ([]T, error)
}

type GetAllRepositoryImpl[T mongodb_domain.BaseDmlModel] struct{}

func GetAllRepository[T mongodb_domain.BaseDmlModel]() IGetAllRepository[T] {
	return &GetAllRepositoryImpl[T]{}
}

func (g *GetAllRepositoryImpl[T]) GetAll(dbName string, collectionName string, createdBy string) ([]T, error) {
	var results []T = make([]T, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongodb_client.GetDb(dbName).Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{"createdBy": createdBy})

	if err != nil {
		log.Printf("❌ [GetAll] error while getting data: %v", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var raw bson.M

		if err := cursor.Decode(&raw); err != nil {
			log.Printf("❌ [GetAll] error while decoding data %v", err)
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
		log.Printf("❌ [GetAll] error in mongodb cursor: %s", err)
		return results, err
	}

	return results, nil
}
