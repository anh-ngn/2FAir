package server

import (
	"log"

	"github.com/bug-breeder/2fair/server/internal/app/handlers"
	"github.com/bug-breeder/2fair/server/internal/pkg/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	// Connect to the database
	err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://2fair.vip"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handlers.SetupRoutes(r)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
