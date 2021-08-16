package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB Name
var dbName = "gonetevent"

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

/* MongoCN : content a MongoDB Connection */
var MongoCN = ConnectDB()

var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

/* ConnectDB : Create a connection to MongoDB and return the connection */
func ConnectDB() MongoInstance {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return MongoInstance{}
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Connect to MongoDB successfully")

	return MongoInstance{Client: client, Db: client.Database(dbName)}
}

/* ConnectionOK: Check connection and return true or false  */
func ConnectionOK() bool {
	err := MongoCN.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}
