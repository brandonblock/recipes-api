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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Recipe is the data model for the recipes our API handles
// swagger:parameters recipes newRecipe
type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

// Store recipes in memory for initial routes
var ctx context.Context
var collection *mongo.Collection

func init() {
	// connect to MongoDB
	ctx = context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// uncomment and run the first time to load from .json to mongo
	// LoadDataToDB()
}

func main() {
	router := gin.Default()

	router.GET("/recipes", ListRecipesHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("recipes/:id", UpdateRecipesHandler)
	router.DELETE("recipes/:id", DeleteRecipeHandler)

	router.Run()
}

// swagger:operation POST /recipes recipes createRecipe
// Creates a new recipe
// ---
// produces:
// - application/json
// responses:
// 200':
// description: Successful operation
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
// 200':
//  description: Successful operation
// '400':
//  description: Invalid input
func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	result := make([]Recipe, 0)

	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		result = append(result, recipe)
	}

	c.JSON(http.StatusOK, result)
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
// '200':
//  description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cur, err := collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	result := make([]Recipe, 0)

	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		result = append(result, recipe)
	}

	c.JSON(http.StatusOK, result)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update and existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//  '200':
//    description: Successful operation
//  '400':
//    description: Invalid input
//  '404':
//    description: Invalid recipe ID
func UpdateRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.UpdateOne(ctx, bson.M{
		"_id": objectID}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: recipe.Name},
			{Key: "instructions", Value: recipe.Instructions},
			{Key: "ingredientes", Value: recipe.Ingredients},
			{Key: "tags", Value: recipe.Tags},
		}}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipes
// Deletes target recipe by id
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
// '200':
// description: Successful operation
// '400':
//  description: Invalid input
//  '404':
//    description: Invalid recipe ID
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Recipe %s has been deleted", id)})
}

// LoadDataToDB is a utility funciont to write sample data from json to the mongo database
func LoadDataToDB() {
	// load recipes from file to memory
	recipes := make([]Recipe, 0)
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
