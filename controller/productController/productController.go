package productController

import (
	"encoding/json"
	"log"
	"net/http"
	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var products []m.Products
	s.DB.Preload("Outlet").Find(&products)
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func Show(c *gin.Context) {
	var product m.Products
	id := c.Param("id")

	if err := s.DB.First(&product, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"detailed_user": product})
}

func Create(c *gin.Context) {
	var product m.Products
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	s.DB.Create(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Data Package Berhasil Ditambahkan"})
}

func Update(c *gin.Context) {
	var product m.Products
	id := c.Param("id")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.DB.Model(&product).Where("id = ?", id).Updates(&product).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat melakukan update package"})
	}
}

func Delete(c *gin.Context) {
	var product m.Products

	var input struct {
		ID json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64()
	if s.DB.Delete(&product, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func ExportToExcel(c *gin.Context) {
	var products []m.Products

	result := s.DB.Debug().Find(&products)
	if result.Error != nil {
		log.Println(result.Error)
	}

	f := excelize.NewFile()

	index, err := f.NewSheet("Sheet1")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	f.SetCellValue("Sheet1", "A1", "No")
	f.SetCellValue("Sheet1", "B1", "ID")
	f.SetCellValue("Sheet1", "C1", "Package Type")
	f.SetCellValue("Sheet1", "D1", "Price")
	f.SetCellValue("Sheet1", "E1", "Signed At")

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	f.SetCellStyle("Sheet1", "A1", "E1", style)

	for i, product := range products {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), product.IDProduct)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), product.Types)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), product.Price)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), product.Outlet.OutletName)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)
	f.SetColWidth("Sheet1", "D", "D", 20)
	f.SetColWidth("Sheet1", "E", "E", 30)

	f.SetActiveSheet(index)

	c.Set("Content-Disposition", "attachment; filename=packages-report.xlsx")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := f.SaveAs("packages-report.xlsx"); err != nil {
		log.Println(err)
		return
	}

	errWrite := f.Write(c.Writer)
	if errWrite != nil {
		log.Println(errWrite)
		return
	}
}

func ExportToPDF(c *gin.Context) {
	var products []m.Products

	result := s.DB.Debug().Find(&products)
	if result.Error != nil {
		log.Println(result.Error)
	}

	f := excelize.NewFile()

	index, err := f.NewSheet("Sheet1")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Package Type")
	f.SetCellValue("Sheet1", "C1", "Price")
	f.SetCellValue("Sheet1", "D1", "Signed At")

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	f.SetCellStyle("Sheet1", "A1", "E1", style)

	for i, product := range products {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), product.IDProduct)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), product.Types)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), product.Price)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), product.Outlet.OutletName)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)
	f.SetColWidth("Sheet1", "D", "D", 30)

	f.SetActiveSheet(index)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	errFont := pdf.AddTTFFont("righteous", "C:/Backup Data/Project/Rich Go/nextlaundry_apis/asset/fonts/Righteous-Regular.ttf")
	if errFont != nil {
		log.Println("failed to add font")
	}
	errFont = pdf.SetFont("righteous", "", 14)
	if errFont != nil {
		log.Println("failed to set font")
	}

	pdf.AddPage()

	r, err := f.GetRows("Sheet1")
	for row, rowCells := range r {
		for _, cell := range rowCells {

			err = pdf.Cell(nil, cell)
			if err != nil {
				log.Println(err)
			}

			pdf.SetX(pdf.GetX() + 100)
		}

		pdf.Br(30)
		pdf.SetX(20)

		if row%20 == 19 {
			pdf.AddPage()
			pdf.SetX(20)
		}

	}

	c.Set("Content-Disposition", "attachment; filename=packages-report.pdf")
	c.Set("Content-Type", "application/pdf")

	if err := pdf.WritePdf("packages-report.pdf"); err != nil {
		log.Println(err)
		return
	}
	errWrite := pdf.Write(c.Writer)
	if errWrite != nil {
		log.Println(errWrite)
		return
	}
}
