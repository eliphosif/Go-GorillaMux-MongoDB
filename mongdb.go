package main

import (
	"context"
	"fmt"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
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

/*
	defer func() {
		if errDB = client.Disconnect(ctx); errDB != nil {
			panic(errDB)
		}
	}()


	podcastcollection := quickStartDB.Collection("podcast")
	episodescollection := quickStartDB.Collection("episode")

	podcastResult, err := podcastcollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "this istest tile"},
		{Key: "suthoer", Value: "Nic niove"},
		{Key: "tags", Value: bson.A{"devel", "priag", "docu"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(podcastResult.InsertedID)

	episodeResult, err := episodescollection.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "podcast", Value: podcastResult.InsertedID},
			{Key: "episode", Value: "episode_1"},
			{Key: "duration", Value: 32},
		},
		bson.D{
			{Key: "podcast", Value: podcastResult.InsertedID},
			{Key: "episode", Value: "episode_2"},
			{Key: "duration", Value: 22},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(episodeResult.InsertedIDs)



	//filter documents
		opts := options.Find()
	opts.SetSort(bson.D{{"age", -1}})
	filterCursor, err := AllCustomers.Find(ctx, bson.D{{"firstname", bson.D{{"$in", bson.A{"Alice", "Bob"}}}}}, opts)
	if err != nil {
		log.Fatal(err)
	}

	var eachCustomer []bson.M

	if err := filterCursor.All(ctx, &eachCustomer); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer filterCursor.Close(ctx)
	for filterCursor.Next(ctx) {
		var eachCustomer bson.M
		if err = filterCursor.Decode(&eachCustomer); err != nil {
			log.Fatal(err)
		}
	}
*/
