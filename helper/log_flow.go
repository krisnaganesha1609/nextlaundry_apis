package helper

import (
	"net/http"
	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"

	"github.com/gin-gonic/gin"
)

func GetLogData(c *gin.Context) {
	var log []m.LogHistory

	s.DB.Find(&log)
	c.JSON(http.StatusOK, gin.H{"log": log})
}
