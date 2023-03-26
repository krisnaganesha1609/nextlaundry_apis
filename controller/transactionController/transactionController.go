package transactionController

import (
	"encoding/json"
	"net/http"

	t "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var trans []t.Transactions

	s.DB.Preload("Placements").Preload("OrderedBy").Preload("InputBy").Find(&trans)
	c.JSON(http.StatusOK, gin.H{"all_transactiondata": trans})
}

func Show(c *gin.Context) {
	var trans t.Transactions
	id := c.Param("id")

	if err := s.DB.Preload("Placement").First(&trans, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data Tidak Ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"detailed_transaction": trans})
}

func Create(c *gin.Context) {
	var trans t.Transactions

	if err := c.ShouldBindJSON(&trans); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	s.DB.Create(&trans)
	c.JSON(http.StatusOK, gin.H{"user": trans})
}

func Update(c *gin.Context) {
	var trans t.Transactions
	id := c.Param("id")

	if err := c.ShouldBindJSON(&trans); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.DB.Model(&trans).Where("id = ?", id).Updates(&trans).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak Dapat Melakukan Update Data"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diperbarui"})
}

func Delete(c *gin.Context) {
	var trans t.Transactions

	var input struct {
		ID json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64()
	if s.DB.Delete(&trans, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus"})
}
