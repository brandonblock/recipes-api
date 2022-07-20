// Recipes API
//
// This is a sample recipes API using Gin.
//
// Schemes http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "recipes-api/handlers"
	models "recipes-api/models"
)

// var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {
	// connect to MongoDB
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if status := redisClient.Ping(); status == nil {
		log.Fatal("unable to connect to redis client")
	}

	// Create route handler
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	// uncomment and run the first time to load from .json to mongo
	// LoadDataToDB(ctx, collection)
}

func main() {
	router := gin.Default()

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.POST("/recipes", recipesHandler.CreateRecipeHandler)
	router.PUT("recipes/:id", recipesHandler.UpdateRecipesHandler)
	router.DELETE("recipes/:id", recipesHandler.DeleteRecipeHandler)

	router.Run()
}

// LoadDataToDB is a utility funciont to write sample data from json to the mongo database
func LoadDataToDB(ctx context.Context, collection *mongo.Collection) {
	// load recipes from file to memory
	recipes := make([]models.Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

	// write in-memory recipes to DB
	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("inserted recipes: ", len(insertManyResult.InsertedIDs))
}
