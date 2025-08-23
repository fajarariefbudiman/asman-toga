package main

import (
	"asman-toga/config"
	"asman-toga/models"
	"asman-toga/router"
	"asman-toga/seeders"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("env file not found, using system env")
	}

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Plant{}, &models.UserPlant{})
	seeders.SeedPlants()
	r := router.SetupRouter()

	r.Run(":8080")
}
