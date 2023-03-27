package dashboardController

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func getMemberCount(c *gin.Context) {
	var count int

	if err := s.DB.Table("tb_member").Select("COUNT(*)").Count(&count); err != nil {
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		
	}

	c.JSON(http.StatusOK, gin.H{"member_count": count})
} 

func getTransactionCount(c *gin.Context) {
	var count int

	if err := s.DB.Table("tb_member").Select("COUNT(*)").Count(&count); err != nil {
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		
	}

	c.JSON(http.StatusOK, gin.H{"member_count": count})
}

func getRevenueCount(c *gin.Context) {
	var count int

	if err := s.DB.Table("tb_member").Select("COUNT(*)").Count(&count); err != nil {
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		
	}

	c.JSON(http.StatusOK, gin.H{"member_count": count})
}