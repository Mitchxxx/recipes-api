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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/gin-contrib/sessions"
	// redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
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

	// Load .env file
  err := godotenv.Load()
  if err != nil {
	log.Fatal("Failed to load .env file")
  }

  ctx = context.Background()
  mongoUri := os.Getenv("MONGO_URI")
  mongoDatabase := os.Getenv("MONGO_DATABASE")
  redisAddress := os.Getenv("REDIS_ADDRESS")
  client, err = mongo.Connect(ctx,
	options.Client().ApplyURI(mongoUri))
  if err = client.Ping(context.TODO(),
	readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collectionRecipes := client.Database(mongoDatabase).Collection("recipes")

	// Code to initialize the Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddress,
		Password: "",
		DB: 0,
	})
	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipeHandler(ctx, collectionRecipes, redisClient)

	collectionUsers := client.Database(mongoDatabase).Collection("users")
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

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	//router.Run()
	
	// Configure and run HTTP server
	server := &http.Server {
		Addr: ":8080",
		Handler: router,
		ReadTimeout:  45 * time.Second,
    	WriteTimeout: 45 * time.Second,
    	IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine
	go func(){
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	log.Println("Server started on :8080")

	// Create Channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV) // catch Ctrl+C and `kill` command
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v/n", err)
	}
	log.Println("Server Stopped")
}