package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Task struct
type Task struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   string
	Desc   string
	Status bool
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Printf("Connecting to mongodb...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Printf("connected!\n")

	log.Println()

	collection := client.Database("tasklist").Collection("tasks")

	res, err := collection.InsertOne(ctx, bson.M{"Name": "Laundry", "Desc": "do the laundry later in the afternoon, 3pm"})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	_, err = collection.InsertOne(ctx, bson.M{"Name": "Study Session: Finals", "Desc": "study for finals"})
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var tasks []Task
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var task Task
		log.Println(cursor)
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		log.Printf("%+v\n", task)
	}

	log.Println()

	var id []byte
	for _, task := range tasks {
		if task.Name == "Laundry" {
			id, err = task.ID.MarshalJSON()
			if err != nil {
				log.Println(err)
			}
		}
	}

	log.Println(string(id))
	str2 := string(id)
	log.Println(len(str2))
	str1 := "5f7b321f0e8d55fd082cb143"
	log.Println(len(str1))
	str := string(id[1:25])
	log.Println(len(str))
	oid, _ := primitive.ObjectIDFromHex(str)
	log.Println(oid)

	log.Println()

	ures, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "Desc", Value: "do the laundry later in the afternoon, 5pm"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ures)

	var taskA Task
	if err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&taskA); err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", taskA)

	log.Println()

	var taskB map[string]interface{}
	if err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&taskB); err != nil {
		log.Println(err)
	}

	log.Printf("%s\n", taskB["_id"])
	log.Printf("%s\n", taskB["Name"])
	log.Printf("%s\n", taskB["Desc"])
	log.Printf("%+v\n", taskB)

	log.Println()

	result, err := collection.DeleteMany(ctx, bson.M{"Name": "Laundry"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	result, err = collection.DeleteOne(ctx, bson.M{"Name": "Study Session: Finals"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

}
