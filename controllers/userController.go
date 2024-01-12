package controllers

import (
	"gin-gorm-tutorial/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

var userBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	if c.Bind(&userBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{Email: userBody.Email, Password: string(hash)}
	r := db.Create(&user)
	if r.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func Login(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	if c.Bind(&userBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	var user *models.User
	db.First(&user, "email = ?", userBody.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, exists := c.Value("user").(*models.User)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
}
