package authController

import (
	"errors"
	"fmt"
	h "nextlaundry_apis/helper"
	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	token, id, err := h.ValidateToken(c)

	if token == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	uid, err := GetUserByID(id.UID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": uid})

}

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AuthHandler(c *gin.Context) {
	var req TokenRequest
	var user m.Users
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	record := s.DB.Preload("Placement").Where("username = ?", req.Username).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	credentialError := user.CheckPasswordHash(req.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials!"})
		c.Abort()
		return
	}

	roles := ""
	switch user.Role {
	case m.Admin:
		roles = "admin"
	case m.Cashier:
		roles = "kasir"
	case m.Owner:
		roles = "owner"
	}

	tokenString, err := h.GenerateJWT(user.IDUser, user.Username, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user.PrepareGive()

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user, "message": "Login Successful"})
}

func LogoutHandler(c *gin.Context) {
	UID := c.MustGet("UID")
	tokenString := h.ExtractToken(c)

	expirationTime := time.Now().Add(24 * time.Hour)
	h.TokenHeap.Push(h.BlacklistedToken{Token: tokenString, ExpirationTime: expirationTime})
	h.BlacklistedTokens[tokenString] = &expirationTime

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d Has Logged Out", UID)})
}

func GetUserByID(uid int) (m.Users, error) {

	var u m.Users

	if err := s.DB.First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil

}

// func RegisterUser(c *gin.Context) {
// 	var user m.Users
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		c.Abort()
// 		return
// 	}

// 	if err := user.HashingPassword(user.Password); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		c.Abort()
// 		return
// 	}

// 	record := s.DB.Create(&user)

// 	if record.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
// 		c.Abort()
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"userId": user.IDUser, "username": user.Username})
// }
