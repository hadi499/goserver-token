package main

import (
	"server-crud/database"
	"server-crud/middleware"

	"server-crud/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Ganti * dengan domain tertentu jika perlu
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika method OPTIONS, langsung response 200 OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	database.ConnectDatabase()

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Static("/uploads", "./uploads")
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.GET("/", controllers.Home)
		authRoutes.POST("/logout", controllers.Logout)
		authRoutes.GET("/products", controllers.GetAllProducts)
		authRoutes.POST("/products", controllers.CreateProduct)
		authRoutes.PUT("/products/:id", controllers.UpdateProduct)
		authRoutes.GET("/products/:id", controllers.GetProductDetail)
		authRoutes.DELETE("/products/:id", controllers.DeleteProduct)
	}
	r.Run(":8080")

}
