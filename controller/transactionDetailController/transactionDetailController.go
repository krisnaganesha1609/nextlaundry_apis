package transactionDetailController

import (
	"encoding/json"
	"net/http"

	t "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var transdet []t.TransactionDetails

	s.DB.Preload("TransactionInfo").Preload("Packages").Find(&transdet)
	c.JSON(http.StatusOK, gin.H{"all_transactiondetailsdata": transdet})
}

func Show(c *gin.Context) {
	var transdet t.TransactionDetails
	id := c.Param("id")

	if err := s.DB.Preload("TransactionInfo").Preload("Packages").First(&transdet, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data Tidak Ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"detailed_transactiondetails": transdet})
}

func Create(c *gin.Context) {
	var transdet []t.TransactionDetails

	if err := c.ShouldBindJSON(&transdet); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	s.DB.Create(&transdet)
	c.JSON(http.StatusOK, gin.H{"message": "Menambah Detail Berhasil"})
}

func Update(c *gin.Context) {
	var transdet t.TransactionDetails
	id := c.Param("id")

	if err := c.ShouldBindJSON(&transdet); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.DB.Model(&transdet).Where("id = ?", id).Updates(&transdet).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak Dapat Melakukan Update Data"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diperbarui"})
}

func Delete(c *gin.Context) {
	var transdet t.TransactionDetails

	var input struct {
		ID json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64()
	if s.DB.Delete(&transdet, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus detail transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus"})
}
