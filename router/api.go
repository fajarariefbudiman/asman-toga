package router

import (
	"asman-toga/controllers"
	"asman-toga/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Asman Toga API running "})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/forgot-password", controllers.ForgotPassword)
		api.POST("/reset-password", controllers.ResetPassword)

		api.GET("/plants", controllers.GetPlants)
		api.GET("/plants/:slug", controllers.GetPlantBySlug)
		api.GET("/userplants", controllers.GetUserPlants)
		api.GET("/userplants/:id", controllers.GetUserPlantByID)
		api.GET("/userplants/:plant_id", controllers.GetUserPlantByPlants)

		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", controllers.GetProfile)
			auth.POST("/logout", controllers.Logout)

			auth.POST("/userplants", controllers.CreateUserPlant)
			auth.PUT("/userplants/:id", controllers.UpdateUserPlant)
			// auth.DELETE("/:id", controllers.DeleteUserPlant)

			auth.PUT("/:id/approve", controllers.ApproveUserPlant)
		}

	}

	return r
}
