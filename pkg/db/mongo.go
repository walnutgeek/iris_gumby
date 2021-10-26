package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userclient *mongo.Client

func ConnectMongoDB() {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second) // context.TODO()

	// user Connection database

	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL")).SetMaxPoolSize(50)

	var err error
	// Connect to MongoDB
	userclient, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = userclient.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to user MongoDB!")

}


//GetMongoDBClient , return mongo client for CRUD operations
func GetMongoDBClient() *mongo.Client {
	return userclient
}
