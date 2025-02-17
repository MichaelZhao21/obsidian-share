package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a custom type for the context key
type contextKey string

const dbContextKey contextKey = "database"

// MongoDB middleware to inject the database into the request context
func WithMongoDB(db *mongo.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dbContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper function to get database from context
func GetDB(r *http.Request) *mongo.Database {
	return r.Context().Value(dbContextKey).(*mongo.Database)
}

func Connect() *mongo.Database {
	// Get the MongoDB URI from the environment
	mongoURI := os.Getenv("MONGODB_URI")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Send ping to confirm successful connections
	var result bson.M
	if err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		log.Fatalf("Error pinging database: %s\n", err.Error())
	}
	fmt.Println("Successfully connected to database!")

	// Get database handle
	db := client.Database("obsidian")
	return db
}
