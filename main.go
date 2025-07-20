//  Recipes API
//
//  This is a simple recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Contact: Mitchel Egboko <megboko@gmail.com>
//  Consumes:
//  - application/json
//
//  Produces
//  - application/json
// swagger:meta

package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mitchxxx/recipes-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var err error
var ctx context.Context
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler


func init(){

	/*recipes = make([]Recipe, 0)

	file, err := os.Open("recipes.json")
	if err != nil {
	   log.Printf("Error opening file: %v", err)
	   return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
	   log.Printf("Error reading file: %v", err)
	   return
	}
	
	if err = json.Unmarshal(fileData, &recipes); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}*/
  ctx = context.Background()
  client, err = mongo.Connect(ctx,
	options.Client().ApplyURI(os.Getenv("MONGO_URI")))
  if err = client.Ping(context.TODO(),
	readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipeHandler(ctx, collection)

	
	/* code to add recipes to mongo database
	var listOfRecipes [] interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
	*/
}

func main () {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	// router.GET("recipes/search", SearchRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.Run()
}