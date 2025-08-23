package controllers

import (
	"asman-toga/config"
	"asman-toga/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Get all UserPlants (bisa dilihat semua user)
func GetUserPlants(c *gin.Context) {
	var userPlants []models.UserPlant

	if err := config.DB.Preload("User").Preload("Plant").Find(&userPlants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil userplants"})
		return
	}

	// Tambah pesan kalau status pending
	var response []gin.H
	for _, up := range userPlants {
		res := gin.H{
			"id":       up.ID,
			"user_id":  up.UserID,
			"plant_id": up.PlantID,
			"address":  up.Address,
			"status":   up.Status,
			"user":     up.User,
			"plant":    up.Plant,
		}
		if up.Status == "pending" {
			res["message"] = "Masih menunggu approval admin"
		}
		response = append(response, res)
	}

	c.JSON(http.StatusOK, response)
}

func GetUserPlantByID(c *gin.Context) {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var userPlant models.UserPlant
	if err := config.DB.Preload("User").Preload("Plant").First(&userPlant, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "UserPlant tidak ditemukan"})
		return
	}

	response := gin.H{
		"id":       userPlant.ID,
		"user_id":  userPlant.UserID,
		"plant_id": userPlant.PlantID,
		"address":  userPlant.Address,
		"status":   userPlant.Status,
		"user":     userPlant.User,
		"plant":    userPlant.Plant,
	}

	if userPlant.Status == "pending" {
		response["message"] = "Masih menunggu approval admin"
	}

	c.JSON(http.StatusOK, response)
}

func ApproveUserPlant(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa approve"})
		return
	}

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var userPlant models.UserPlant
	if err := config.DB.First(&userPlant, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "UserPlant tidak ditemukan"})
		return
	}

	userPlant.Status = "approved"

	if err := config.DB.Save(&userPlant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve userplant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userplant": userPlant})
}

func CreateUserPlant(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		userIDStr, strOk := userIDVal.(string)
		if !strOk {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
			return
		}

		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
			return
		}
	}

	var input struct {
		PlantID   int     `json:"plant_id" binding:"required"`
		Address   string  `json:"address" binding:"required"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Notes     string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plant models.Plant
	if err := config.DB.First(&plant, input.PlantID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plant not found"})
		return
	}

	userPlant := models.UserPlant{
		ID:        uuid.New(),
		UserID:    userID,
		PlantID:   input.PlantID,
		Address:   input.Address,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Notes:     input.Notes,
		Status:    "pending",
	}

	if err := config.DB.Create(&userPlant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal tambah userplant", "details": err.Error()})
		return
	}

	config.DB.Preload("Plant").Preload("User").First(&userPlant, userPlant.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":   "UserPlant berhasil dibuat, menunggu persetujuan admin",
		"userplant": userPlant,
	})
}

func UpdateUserPlant(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		userIDStr, strOk := userIDVal.(string)
		if !strOk {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
			return
		}

		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
			return
		}
	}

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var userPlant models.UserPlant
	if err := config.DB.First(&userPlant, "id = ? AND user_id = ?", uid, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "UserPlant tidak ditemukan"})
		return
	}

	var input struct {
		Address   string  `json:"address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Notes     string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Address != "" {
		userPlant.Address = input.Address
	}
	if input.Latitude != 0 {
		userPlant.Latitude = input.Latitude
	}
	if input.Longitude != 0 {
		userPlant.Longitude = input.Longitude
	}
	if input.Notes != "" {
		userPlant.Notes = input.Notes
	}

	userPlant.Status = "pending"

	if err := config.DB.Save(&userPlant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update userplant"})
		return
	}
	config.DB.Preload("Plant").Preload("User").First(&userPlant, userPlant.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "UserPlant berhasil diperbarui, menunggu persetujuan admin",
		"userplant": userPlant,
	})
}
