package controllers

import (
	"asman-toga/config"
	"asman-toga/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllBanjar(c *gin.Context) {
	var banjars []models.Banjar
	if err := config.DB.Find(&banjars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data tanaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"banjars": banjars})
}
