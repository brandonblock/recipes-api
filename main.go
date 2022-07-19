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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "recipes-api/handlers"
	models "recipes-api/models"
)

// Store recipes in memory for initial routes
var ctx context.Context
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {
	// connect to MongoDB
	ctx = context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// Create route handler
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)

	// uncomment and run the first time to load from .json to mongo
	// LoadDataToDB()
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
func LoadDataToDB() {
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
