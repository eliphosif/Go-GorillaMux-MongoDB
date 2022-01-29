package main

import (
	"context"
	"fmt"
	"log"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

var AllCustomers *mongo.Collection
var ctx = context.Background()

func initlizeMongoConnection() *mongo.Collection {

	MongoDBURI := "mongodb+srv://eliphosif:eliphosif@cluster0.0imgv.mongodb.net/test?authSource=admin&replicaSet=atlas-fydu9m-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"

	//defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI))

	golangMongoDB := client.Database("GolangMongo")
	AllCustomers := golangMongoDB.Collection("AllCustomers")
	return AllCustomers

}

func insertCustomerDoc(AllCustomers *mongo.Collection, cust Customer, ctx context.Context) (*mongo.InsertOneResult, error) {
	//insert document in to MongoDB collection
	res, err := AllCustomers.InsertOne(ctx, cust)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}
