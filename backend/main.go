package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Get the secret from the environment
	hmacSecret := os.Getenv("SUPABASE_JWT_SECRET")

	// Prevent the server from starting if the secret is not set
	if hmacSecret == "" {
		log.Fatal("Please set the SUPABASE_JWT_SECRET environment variable")
	}

	// Create a new router
	router := gin.New()

	// Enable CORS for all origins. This is not recommended for production usage.
	// Use a whitelist of allowed origins instead.
	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}

	router.Use(cors.New(corsConfig))

	// The only route we have is /secret and it is protected by the authMiddleware.
	router.POST("/secret", authMiddleware(hmacSecret), secretRouteHandler())

	// Run the server
	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

var emailCtxKey = "email"

func authMiddleware(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Validate token
		// convert strign to a byte array
		email, err := parseJWTToken(token, []byte(hmacSecret))

		if err != nil {
			log.Printf("Error parsing token: %s", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		log.Printf("Received request from %s", email)

		// Save the email in the context to use later in the handler
		ctx := context.WithValue(c, emailCtxKey, email)
		c.Request = c.Request.WithContext(ctx)

		// Authenticated. Continue (call next handler)
		c.Next()
	}
}

func secretRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the email from the context
		email := c.GetString(emailCtxKey)

		// Return the secret message
		c.JSON(200, gin.H{
			"message": "our hidden value for the user " + email,
		})
	}
}

// List of claims that we want to parse from the JWT token.
// The RegisteredClaims struct contains the standard claims.
// See https://pkg.go.dev/github.com/golang-jwt/jwt/v5#RegisteredClaims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// This function parses the JWT token and returns the email claim
func parseJWTToken(token string, hmacSecret []byte) (email string, err error) {
	// Parse the token and validate the signature
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	// Check if the token is valid
	if err != nil {
		return "", fmt.Errorf("error validating token: %v", err)
	} else if claims, ok := t.Claims.(*Claims); ok {
		return claims.Email, nil
	}

	return "", fmt.Errorf("error parsing token: %v", err)
}
