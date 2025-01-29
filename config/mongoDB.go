package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func DbConnect() {
	var err error
	uri := "mongodb+srv://deep82500:deep82500@deep.jqe1i.mongodb.net/?retryWrites=true&w=majority&appName=deep"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbLocation := options.Client().ApplyURI(uri)
	Client, err = mongo.Connect(ctx, dbLocation)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("failed to connect with db", err.Error())
	}
	log.Println("successfully connected to the mongoDB")

}

func Getcollection() *mongo.Collection {

	Collection = Client.Database("user_info(PB)").Collection("info")
	return Collection

}
