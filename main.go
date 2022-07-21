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
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "recipes-api/handlers"
	models "recipes-api/models"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler

func init() {
	// connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collectionRecipes := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})
	if status := redisClient.Ping(); status == nil {
		log.Fatal("unable to connect to redis client")
	}

	// Create route handler
	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)

	// Create auth handler
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

	// uncomment and run to load test data to db
	// LoadRecipDataToDB(ctx, collectionRecipes)
	// LoadUsersDataToDB(ctx, collectionUsers)
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", os.Getenv("REDIS_URI"), "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))

	// unsecured endpoints
	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signout", authHandler.SignOutHandler)

	// secured endpoints
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
		authorized.POST("/recipes", recipesHandler.CreateRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipesHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)

		authorized.POST("/refresh", authHandler.RefreshHandler)
	}

	router.Run()
}

// LoadDataToDB is a utility funciont to write sample data from json to the mongo database
func LoadRecipeDataToDB(ctx context.Context, collection *mongo.Collection) {
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

func LoadUsersDataToDB(ctx context.Context, collection *mongo.Collection) {
	users := map[string]string{
		"admin":  "fCRmh4Q2J7Rseqkz",
		"bblock": "123password",
	}
	h := sha256.New()
	for uname, pword := range users {
		hashed := string(h.Sum([]byte(pword)))
		collection.InsertOne(ctx, bson.M{
			"username": uname,
			"password": hashed,
		})
	}
}
