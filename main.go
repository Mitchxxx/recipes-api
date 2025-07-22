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
	"github.com/go-redis/redis"
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
var authHandler *handlers.AuthHandler


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
	collectionRecipes := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// Code to initialize the Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipeHandler(ctx, collectionRecipes, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
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

func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context) {
		if c.GetHeader("X-API-KEY") !=
		os.Getenv("X_API_KEY") {
			c.AbortWithStatus(401)
		}
	}
}


func main () {
	router := gin.Default()

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SigninHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	router.Run()
}