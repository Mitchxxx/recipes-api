package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx context.Context
}

func NewAuthHandler (ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx: ctx,
	}
}


func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		auth0Domain := os.Getenv("AUTH_DOMAIN")
		apiIdentifier := os.Getenv("AUTH0_API_IDENTIFIER")

		if auth0Domain == "" || apiIdentifier == "" {
			log.Fatal("AUTH_DOMAIN or AUTH)_API_IDENTIFIER env variable not set")
		}

		auth0Domain = "https://" + auth0Domain + "/"

		log.Println("Auth0 Domain:", auth0Domain)
		log.Println("API Identifier", apiIdentifier)
		client := auth0.NewJWKClient(auth0.JWKClientOptions{
			URI: auth0Domain + ".well-known/jwks.json",
			}, nil)
		configuration := auth0.NewConfiguration(
			client,
			[]string{apiIdentifier},
			auth0Domain,
			jose.RS256,
		)

		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token validation failed",
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
