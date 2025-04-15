package main

import (
	"front-back_6/handlers"
	"front-back_6/middleware"
	"front-back_6/models"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func DbInit() *gorm.DB {
	db, err := models.Setup()
	if err != nil {
		log.Println("Can't connect to database")
	}

	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default() // Создание экземпляра сервера

	db := DbInit() // Инициализация базы данных

	server := handlers.NewServer(db) // Создание экземпляра сервера

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := r.Group("/api")

	router.POST("/register", server.Register)
	router.POST("/login", server.Login)

	user := r.Group("/user")
	user.Use(middleware.JWTMiddleware())
	user.GET("/profile", server.GetCurrentUser)
	user.POST("/logout", server.LogOut)
	user.GET("/data", server.GetData)

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file") // Загрузка файла .env
	}
	port := os.Getenv("PORT") // Получение порта из файла .env

	r := SetupRouter() // Создание экземпляра сервера

	log.Fatal(r.Run(":" + port)) // Запуск сервера
}
