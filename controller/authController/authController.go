package authController

import (
	"errors"
	h "nextlaundry_apis/helper"
	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func SecondValidate(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "token valid", "data": uid})

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

	s.DB.Exec("INSERT INTO log_history VALUES(NULL, CONCAT('New login activity by: ', ?), 'login', NOW(), NULL);", user.IDUser)
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user, "message": "Logged In Successfully"})
}

func LogoutHandler(c *gin.Context) {
	tokenString := h.ExtractToken(c)

	expirationTime := time.Now().Add(24 * time.Hour)
	h.TokenHeap.Push(h.BlacklistedToken{Token: tokenString, ExpirationTime: expirationTime})
	h.BlacklistedTokens[tokenString] = &expirationTime

	c.JSON(http.StatusOK, gin.H{"message": "Logged Out Successfully"})
}

func GetUserByID(uid int) (m.Users, error) {

	var u m.Users

	if err := s.DB.Preload("Placement").First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil

}
