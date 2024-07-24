package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bug-breeder/2fair/server/configs"
	"github.com/bug-breeder/2fair/server/internal/app/models"
	"github.com/bug-breeder/2fair/server/internal/pkg/auth"
	"github.com/bug-breeder/2fair/server/internal/pkg/db"
	"github.com/gorilla/sessions"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine) {
	// Ensure SESSION_SECRET is set
	gothic.Store = sessions.NewCookieStore([]byte(configs.GetEnv("SESSION_SECRET")))

	redirectURI := os.Getenv("REDIRECT_URI")

	goth.UseProviders(
		google.New(configs.GetEnv("GOOGLE_CLIENT_ID"), configs.GetEnv("GOOGLE_CLIENT_SECRET"), redirectURI+"/auth/google/callback"),
		microsoftonline.New(configs.GetEnv("MICROSOFT_CLIENT_ID"), configs.GetEnv("MICROSOFT_CLIENT_SECRET"), redirectURI+"/auth/microsoft/callback"),
		// apple.New(configs.GetEnv("APPLE_CLIENT_ID"), configs.GetEnv("APPLE_CLIENT_SECRET"), "http://localhost:8080/auth/apple/callback"),
	)

	router.GET("/auth/:provider", authHandler)
	router.GET("/auth/:provider/callback", authCallback)
	router.POST("/auth/refresh", RefreshToken)

	protected := router.Group("/")
	protected.Use(auth.Authenticate())
	protected.POST("/otps", AddOTP)
	protected.PUT("/otps/:otpID/inactivate", InactivateOTP)
	protected.PUT("/otps/:otpID", EditOTP)
	protected.GET("/otps", ListOTPs)
	protected.GET("/otps/codes", GenerateOTPCodes)
}

func authHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider is required"})
		return
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func authCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Failed to complete user auth: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete user auth"})
		return
	}

	// Store user information in MongoDB
	usersCollection := db.GetCollection("users")

	// Check if a user with the same email already exists
	existingUser := models.User{}
	err = usersCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to check existing user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
		return
	}

	var accessToken, refreshToken string

	if err == mongo.ErrNoDocuments {
		// No existing user found, create a new one
		newUser := models.User{
			Name:       user.Name,
			Email:      user.Email,
			Provider:   user.Provider,
			ProviderID: user.UserID,
			CreatedAt:  time.Now(),
			OTPs:       []models.OTP{},
		}

		result, err := usersCollection.InsertOne(context.Background(), newUser)
		if err != nil {
			log.Printf("Failed to create new user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user"})
			return
		}

		userID := result.InsertedID.(primitive.ObjectID).Hex()
		accessToken, err = auth.GenerateAccessToken(userID, newUser.Email)
		if err != nil {
			log.Printf("Failed to generate access token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err = auth.GenerateRefreshToken(userID, newUser.Email)
		if err != nil {
			log.Printf("Failed to generate refresh token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

	} else {
		// Existing user found, return the user info without updating
		userID := existingUser.ID.Hex()
		accessToken, err = auth.GenerateAccessToken(userID, existingUser.Email)
		if err != nil {
			log.Printf("Failed to generate access token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err = auth.GenerateRefreshToken(userID, existingUser.Email)
		if err != nil {
			log.Printf("Failed to generate refresh token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}
	}

	// Set tokens in cookies
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	c.JSON(http.StatusOK, gin.H{
		// "user":         user,
		"access_token": accessToken,
	})
}