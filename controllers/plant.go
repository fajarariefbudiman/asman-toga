package controllers

import (
	"net/http"

	"asman-toga/config"
	"asman-toga/models"

	"github.com/gin-gonic/gin"
)

func GetPlants(c *gin.Context) {
	var plants []models.Plant
	if err := config.DB.Find(&plants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data tanaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plants": plants})
}

func GetPlantBySlug(c *gin.Context) {
	slug := c.Param("slug")
	// plantID, err := uuid.Parse(idParam)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Data tanaman tidak valid"})
	// 	return
	// }

	var plant models.Plant
	if err := config.DB.First(&plant, "slug = ?", slug).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tanaman tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plant": plant})
}
