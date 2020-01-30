package mongoClient

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"
	"fmt"
)

type MongoClient struct {
	MongoClient *mongo.Client
}

func CreateConnectedClient() MongoClient {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://:27017/"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic("Cannot connect to mongo.")
	}
	fmt.Println("CONNECTED tO MONGO")
	return MongoClient{ MongoClient: client }
}
