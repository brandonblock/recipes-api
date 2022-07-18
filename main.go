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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

//Recipe is the data model for the recipes our API handles
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// Store recipes in memory for initial routes
var recipes []Recipe

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

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

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
	c.JSON(http.StatusOK, recipes)
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
	result := make([]Recipe, 0)

	for _, r := range recipes {
		found := false
		for _, t := range r.Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			result = append(result, r)
		}
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

	//nasty iteration because we're storing in a huge array
	index := -1
	for i, r := range recipes {
		if r.ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe Not Found"})
		return
	}
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
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

	index := -1
	for i, r := range recipes {
		if r.ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe Not Found"})
		return
	}
	// amazingly dumb array manipulation to keep this all in memory
	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Recipe %s has been deleted", id)})
}

// while recipes are held in memory, we'll load them from a JSON file at startup
func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
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
