package seeders

import (
	"fmt"

	"asman-toga/config"
	"asman-toga/models"
)

func SeedPlants() {
	plants := []models.Plant{
		{PlantName: "Lidah Buaya", Slug: "lidah-buaya"},
		{PlantName: "Jahe", Slug: "jahe"},
		{PlantName: "Jahe Hitam", Slug: "jahe-hitam"},
		{PlantName: "Dadap", Slug: "dadap"},
		{PlantName: "Sambiloto", Slug: "sambiloto"},
		{PlantName: "Kitolod", Slug: "kitolod"},
		{PlantName: "Kunyit", Slug: "kunyit"},
		{PlantName: "Kunyit Putih", Slug: "kunyit-putih"},
		{PlantName: "Serai", Slug: "serai"},
		{PlantName: "Daun Sirih", Slug: "daun-sirih"},
	}

	for _, plant := range plants {
		var existing models.Plant
		if err := config.DB.Where("plant_name = ?", plant.PlantName).First(&existing).Error; err == nil {
			continue
		}
		config.DB.Create(&plant)
		fmt.Println("✅ Seeded plant:", plant.PlantName)
	}
}
