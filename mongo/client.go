package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
	databaseName   string
)

var (
	ErrNoDocuments = mongo.ErrNoDocuments
)

// LoadConfig loads Mongo config from environment
func LoadConfig() (uri string, db string, maxPool uint64, minPool uint64, err error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	uri = os.Getenv("MONGO_URI")
	db = os.Getenv("MONGO_DB")
	if uri == "" || db == "" {
		return "", "", 0, 0, fmt.Errorf("MONGO_URI or MONGO_DB not set")
	}

	maxPoolStr := os.Getenv("MONGO_MAXPOOLSIZE")
	minPoolStr := os.Getenv("MONGO_MINPOOLSIZE")

	maxPool = 100
	minPool = 0

	if maxPoolStr != "" {
		if v, err := strconv.ParseUint(maxPoolStr, 10, 64); err == nil {
			maxPool = v
		}
	}
	if minPoolStr != "" {
		if v, err := strconv.ParseUint(minPoolStr, 10, 64); err == nil {
			minPool = v
		}
	}

	return uri, db, maxPool, minPool, nil
}

// GetClient returns a singleton Mongo client
func GetClient() (*mongo.Client, string, error) {
	var err error
	clientOnce.Do(func() {
		uri, db, maxPool, minPool, e := LoadConfig()

		databaseName = db
		if e != nil {
			err = e
			return
		}
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetMaxPoolSize(maxPool).
			SetMinPoolSize(minPool)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)

		if err != nil {
			panic(err)
		}
		clientInstance = client
		// Send a ping to confirm a successful connection
		if err := clientInstance.Database(db).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
			panic(err)
		}
		fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	})

	return clientInstance, databaseName, err
}

// GetCollection helper
func GetCollection(collectionName string) (*mongo.Collection, error) {
	client, db, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Database(db).Collection(collectionName), nil
}

// Disconnect safely closes the Mongo client
func Disconnect() error {
	if clientInstance == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return clientInstance.Disconnect(ctx)
}
