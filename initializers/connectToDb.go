package initializers

import (
	"context"
	"log"
	"os"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var GroupCollection *mongo.Collection
var SessionCollection *mongo.Collection
var UserTodoCollection *mongo.Collection
var GroupTodoCollection *mongo.Collection

var MongoClient *mongo.Client

func ConnectToDB() {

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	MongoClient = client

	if err != nil {
		log.Println(err)
		return
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Err(); err != nil {
		log.Println(err)
		return
	}

	UserCollection = client.Database("Users").Collection("user_data")
	UserTodoCollection = client.Database("Users").Collection("user_todos")
	GroupCollection = client.Database("Groups").Collection("group_data")
	GroupTodoCollection = client.Database("Groups").Collection("group_todos")
	SessionCollection = client.Database("Session").Collection("session_data")

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
