package main

import (
	"api/service/controllers"
	"api/service/db"
	"api/service/middleware"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.DbOpen()

	router := gin.Default()

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000" // fallback para dev
	}

	// Middleware CORS configurado corretamente
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // ← necessário para cookies
		MaxAge:           12 * time.Hour, // cache do preflight
	}))

	api := router.Group("/api")
	api.Use(middleware.LoggedMiddleware)

	{
		api.GET("/agendamentos/:id", controllers.GetAgendamentos)
		api.GET("/saidas", controllers.GetSaidas)
		api.POST("/saidas", controllers.PostSaidas)
		api.DELETE("/saidas/:id", controllers.DeleteSaidas)
		api.PUT("/saidas/:id", controllers.PutSaidas)
	}

	router.POST("/login", controllers.LoginController)

	router.Run("localhost:3000")

}
