package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchxxx/recipes-api/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token string `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler (ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx: ctx,
	}
}

func (handler *AuthHandler) SigninHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	storedUser := models.User{}
	err := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
	}).Decode(&storedUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	// Compare passwords using bcrypt
	if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}
	// Generate a unique session token
	sessionToken := xid.New().String()

	// Create and save session
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	if err := session.Save(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save session",
        })
        return
    }
	c.JSON(http.StatusOK, gin.H{
		"message": "User signed in",
		"username": storedUser.Username,
	})

}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

func (handler *AuthHandler) RefreshHandler (c *gin.Context){
	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")

	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid session cookie",
		})
		return
	}
	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "New session issued",
	})

}

func (handler *AuthHandler) SignOutHandler (c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out ...",
	})
}