package dao

import (
	"context"
	"log"
	"math/rand"

	"os"
	"time"

	"github.com/St0iK/go-quote-parser/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection

const (
	DBNAME     = "quotes-parser"
	COLLECTION = "quotes"
)

// Connect ...
func Connect() {
	log.Println("Initialising MongoDB connection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_DB_URL")))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(DBNAME).Collection(COLLECTION)
	log.Println("Connected to MongoDB!")

}

// something
func GetRandomQuote() model.Quote {

	rand.Seed(time.Now().UnixNano())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	c, _ := collection.CountDocuments(ctx, bson.D{}, nil)

	random := rand.Int63n(c-1) + 1

	findOptions := options.Find()
	findOptions.SetLimit(1)
	findOptions.SetSkip(random)

	cur, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem model.Quote
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		return elem
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return model.Quote{}
}
