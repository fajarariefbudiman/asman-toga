package controllers

import (
	"asman-toga/config"
	"asman-toga/models"
	"asman-toga/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	// otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	// reset := models.ResetPassword{
	// 	ID:        uuid.New(),
	// 	Email:     req.Email,
	// 	OTP:       otp,
	// 	ExpiresAt: time.Now().Add(10 * time.Minute),
	// }
	// config.DB.Create(&reset)
	var reset models.ResetPassword

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	if err := config.DB.Where("email = ?", req.Email).First(&reset).Error; err == nil {
		// update OTP lama
		reset.OTP = otp
		reset.ExpiresAt = time.Now().Add(10 * time.Minute)
		config.DB.Save(&reset)
	} else {
		// kalau belum ada, bikin baru
		reset = models.ResetPassword{
			ID:        uuid.New(),
			Email:     req.Email,
			OTP:       otp,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}
		config.DB.Create(&reset)
	}

	go utils.SendEmail(req.Email, "Reset Password OTP", "Kode OTP kamu: "+otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sudah dikirim ke email"})
}

func ResetPassword(c *gin.Context) {
	var req struct {
		Email           string `json:"email"`
		OTP             string `json:"otp"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak sama"})
		return
	}

	var reset models.ResetPassword
	if err := config.DB.Where("email = ? AND otp = ?", req.Email, req.OTP).
		Order("created_at desc").First(&reset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP salah"})
		return
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP sudah kadaluarsa"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	config.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("password", string(hashed))

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset"})
}
