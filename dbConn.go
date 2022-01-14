package main

import (
	"context"
	"fmt"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//const database, collection = "meteringdb", "metrics"

var mongoConn = MongoConn()

func MongoConn() (client *mongo.Client) {
	credential := options.Credential{
		AuthSource: "meteringdb",
		Username:   "user",
		Password:   "pwd",
	}
	godotenv.Load(".env")
	mongoUrl := os.Getenv("MONGO_URL")

	clientOptions := options.Client().ApplyURI(mongoUrl).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB Connection Made")
	return client
}

func GetCollection() (col *mongo.Collection) {
	database := os.Getenv("MONGO_DBNAME")
	collection := os.Getenv("MONGO_COLLECTION")

	metrics := mongoConn.Database(database).Collection(collection)

	return metrics
}

func InfluxConn() influxdb2.Client {
	godotenv.Load(".env")
	influxUrl := os.Getenv("INFLUX_URL")
	influxToken := os.Getenv("INFLUX_TOKEN")

	client := influxdb2.NewClient(influxUrl, influxToken)

	return client
}
