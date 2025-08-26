package seeders

import (
	"asman-toga/config"
	"asman-toga/models"
	"fmt"
)

func SeedBanjar() {
	BanjarArray := []models.Banjar{
		{Name: "Patus", Slug: "patus"},
		{Name: "Tengah", Slug: "tengah"},
		{Name: "Buayang", Slug: "buayang"},
		{Name: "Bandung", Slug: "bandung"},
		{Name: "Nyamping", Slug: "nyamping"},
		{Name: "Babung", Slug: "babung"},
		{Name: "Kebon", Slug: "kebon"},
	}

	for _, Banjar := range BanjarArray {
		var existing models.Banjar
		if err := config.DB.Where("name = ?", Banjar.Name).First(&existing).Error; err == nil {
			continue
		}
		config.DB.Create(&Banjar)
		fmt.Println("✅ Seeded Banjar:", Banjar.Name)
	}
}
