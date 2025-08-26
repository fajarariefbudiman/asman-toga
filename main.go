package main

import (
	"asman-toga/config"
	"asman-toga/models"
	"asman-toga/router"
	"asman-toga/seeders"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("env file not found, using system env")
	}

	config.ConnectDatabase()

	// config.DB.AutoMigrate(&models.User{}, &models.Plant{}, &models.UserPlant{}, &models.Banjar{}, &models.ResetPassword{})
	error := config.DB.AutoMigrate(
		&models.User{},
		&models.Plant{},
		&models.UserPlant{},
		&models.Banjar{},
		&models.ResetPassword{},
	)
	if error != nil {
		log.Fatal("AutoMigrate gagal:", error)
	}

	seeders.SeedPlants()
	seeders.SeedBanjar()
	r := router.SetupRouter()

	r.Run(":8080")
}
