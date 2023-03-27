package memberController

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var members []m.Members
	s.DB.Find(&members)
	c.JSON(http.StatusOK, gin.H{"members": members})
}

func GetByMonth(c *gin.Context) {
	var members []m.Members
	var date m.Date

	if err := c.ShouldBindJSON(&date); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	s.DB.Where(
		"created_at BETWEEN ? AND ?", date.StartDate, date.EndDate,
	).Find(&members)

	c.JSON(http.StatusOK, gin.H{
		"label": "member",
		"data":  members,
	})
}

func Show(c *gin.Context) {
	var member m.Members
	id := c.Param("id")

	if err := s.DB.First(&member, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data Tidak Ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"detailed_member": member})
}

func Create(c *gin.Context) {
	var member m.Members

	if err := c.ShouldBindJSON(&member); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	s.DB.Create(&member)
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil Menambahkan Member"})
}

func Update(c *gin.Context) {
	var member m.Members
	id := c.Param("id")

	if err := c.ShouldBindJSON(&member); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.DB.Model(&member).Where("id = ?", id).Updates(&member).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak Dapat Melakukan Update Member"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diperbarui"})
}

func Delete(c *gin.Context) {
	var member m.Members

	var input struct {
		ID json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64()
	if s.DB.Delete(&member, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus Member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus"})
}

func ExportToExcel(c *gin.Context) {
	var members []m.Members

	result := s.DB.Debug().Find(&members)
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
	f.SetCellValue("Sheet1", "C1", "Name")
	f.SetCellValue("Sheet1", "D1", "Address")
	f.SetCellValue("Sheet1", "E1", "Gender")
	f.SetCellValue("Sheet1", "F1", "Phone")

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

	f.SetCellStyle("Sheet1", "A1", "F1", style)

	for i, member := range members {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), member.IDMember)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), member.MemberName)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), member.MemberAddress)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), member.Gender)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), member.MemberPhone)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)
	f.SetColWidth("Sheet1", "D", "D", 20)
	f.SetColWidth("Sheet1", "E", "E", 30)
	f.SetColWidth("Sheet1", "F", "F", 20)

	f.SetActiveSheet(index)

	c.Set("Content-Disposition", "attachment; filename=members-report.xlsx")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := f.SaveAs("members-report.xlsx"); err != nil {
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
	var members []m.Members

	result := s.DB.Debug().Find(&members)
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
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "Username")

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

	for i, member := range members {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), member.IDMember)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), member.MemberName)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), member.MemberPhone)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)

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

	c.Set("Content-Disposition", "attachment; filename=members-report.pdf")
	c.Set("Content-Type", "application/pdf")

	if err := pdf.WritePdf("members-report.pdf"); err != nil {
		log.Println(err)
		return
	}
	errWrite := pdf.Write(c.Writer)
	if errWrite != nil {
		log.Println(errWrite)
		return
	}
}
