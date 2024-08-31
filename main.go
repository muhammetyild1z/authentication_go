package main

import (
	"authentication_go/database"
	"authentication_go/middleware"
	"authentication_go/repository"
	"authentication_go/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Veritabanı başlatma
	database.InitDatabase()

	// Repository ve servisleri oluşturma
	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)

	// Gin başlatma
	r := gin.Default()

	// Kayıt olma endpoint'i
	r.POST("/register", func(c *gin.Context) {
		var json struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := userService.RegisterUser(json.Username, json.Password, json.Role); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User registered successfully"})
	})

	// Login endpoint'i
	r.POST("/login", func(c *gin.Context) {
		var json struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		token, err := userService.AuthenticateUser(json.Username, json.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	})

	// Sadece admin yetkilendirilmiş endpoint
	r.GET("/admin", middleware.AuthMiddleware("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome, admin!"})
	})

	r.Run(":8080")
}
