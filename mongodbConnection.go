package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://cherindranath:cherry2580@cluster0.s1pf5.mongodb.net/publish_msg.msg?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	//data := bson.D{{Key: "name", Value: "this is a test msg"}}
	//insertIntoDB(client, data)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func insertIntoDB(client *mongo.Client, data bson.D) {
	collection := client.Database("publish_msg").Collection("msg")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(id)
}
