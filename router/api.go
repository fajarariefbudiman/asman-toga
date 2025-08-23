package router

import (
	"asman-toga/controllers"
	"asman-toga/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Asman Toga API running 🚀"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		api.GET("/plants", controllers.GetPlants)
		api.GET("/plants/:slug", controllers.GetPlantBySlug)
		api.GET("/userplants", controllers.GetUserPlants)
		api.GET("/userplants/:id", controllers.GetUserPlantByID)

		auth := api.Group("/userplants")
		auth.Use(middleware.JWTAuth())
		{
			auth.POST("/", controllers.CreateUserPlant)
			auth.PUT("/:id", controllers.UpdateUserPlant)
			// auth.DELETE("/:id", controllers.DeleteUserPlant)

			auth.PUT("/:id/approve", controllers.ApproveUserPlant)
		}

	}

	return r
}
