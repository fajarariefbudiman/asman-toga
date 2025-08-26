package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"asman-toga/config"
	"asman-toga/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var validate = validator.New()

func translateError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " wajib diisi"
	case "min":
		return fe.Field() + " minimal " + fe.Param() + " karakter"
	case "max":
		return fe.Field() + " maksimal " + fe.Param() + " karakter"
	case "email":
		return "Format email tidak valid"
	case "oneof":
		return fe.Field() + " hanya boleh salah satu dari: " + fe.Param()
	default:
		return fe.Field() + " tidak valid"
	}
}

func Register(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required,min=3,max=100"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		BanjarId        int    `json:"banjar_id" binding:"required"`
		Role            string `json:"role" binding:"omitempty,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errMsgs := []string{}
			for _, fe := range verrs {
				errMsgs = append(errMsgs, translateError(fe))
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"errors":  errMsgs,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Password dan Konfirmasi Password tidak sama",
		})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    strings.ToLower(req.Email),
		Password: string(hashed),
		BanjarID: req.BanjarId,
		Role:     req.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email sudah terdaftar",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errMsgs := []string{}
			for _, fe := range verrs {
				errMsgs = append(errMsgs, translateError(fe))
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"errors":  errMsgs,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email atau password salah",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email atau password salah",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

var blacklistedTokens = make(map[string]time.Time)

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	blacklistedTokens[tokenString] = time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func IsBlacklisted(token string) bool {
	exp, exists := blacklistedTokens[token]
	if !exists {
		return false
	}
	if time.Now().After(exp) {
		delete(blacklistedTokens, token)
		return false
	}
	return true
}
