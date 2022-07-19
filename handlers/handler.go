package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"recipes-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation POST /recipes recipes createRecipe
// Creates a new recipe
// ---
// produces:
// - application/json
// responses:
// 200':
// description: Successful operation
func (handler *RecipesHandler) CreateRecipeHandler(c *gin.Context) {
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := handler.collection.InsertOne(handler.ctx, recipe)
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
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	result := make([]models.Recipe, 0)

	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		result = append(result, recipe)
	}

	c.JSON(http.StatusOK, result)
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Returns list of recipes, search by tag
// ---
// produces:
// - application/json
// responses:
// '200':
//  description: Successful operation
func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cur, err := handler.collection.Find(handler.ctx, bson.M{"tags": tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	result := make([]models.Recipe, 0)

	for cur.Next(handler.ctx) {
		var recipe models.Recipe
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
func (handler *RecipesHandler) UpdateRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
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
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectID})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Recipe %s has been deleted", id)})
}
