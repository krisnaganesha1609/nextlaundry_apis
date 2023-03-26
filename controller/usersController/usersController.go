package usersController

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

func Index(c *gin.Context) {
	var users []m.Users

	s.DB.Preload("Placement").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func Show(c *gin.Context) {
	var user m.Users
	id := c.Param("id")

	if err := s.DB.Preload("Placement").First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data Tidak Ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"detailed_user": user})
}

func Create(c *gin.Context) {
	var user m.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	newUser := m.Users{
		Fullname:  user.Fullname,
		Username:  user.Username,
		Role:      user.Role,
		Id_outlet: user.Id_outlet,
	}

	hashed, err := user.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newUser.Password = hashed

	errCreateUser := s.DB.Create(&newUser).Error
	if errCreateUser != nil {
		errCreateUser := strings.Split(errCreateUser.Error(), ":")[0]
		log.Println(errCreateUser)
		if errCreateUser == "Error 1062 (23000)" {
			c.AbortWithStatusJSON(400, gin.H{"message": "Duplicate Entry For This Data"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil Menambahkan User Baru"})
}

func Update(c *gin.Context) {
	var user m.Users
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.DB.Model(&user).Where("id = ?", id).Updates(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak Dapat Melakukan Update User"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diperbarui"})
}

func Delete(c *gin.Context) {
	var user m.Users

	var input struct {
		ID json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64()
	if s.DB.Delete(&user, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus"})
}

func ExportToExcel(c *gin.Context) {
	var users []m.Users

	result := s.DB.Debug().Find(&users)
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
	f.SetCellValue("Sheet1", "D1", "Username")
	f.SetCellValue("Sheet1", "E1", "Placement")
	f.SetCellValue("Sheet1", "F1", "Roled As")

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

	for i, user := range users {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), user.IDUser)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), user.Fullname)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), user.Username)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), user.Placement.OutletName)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), user.Role)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)
	f.SetColWidth("Sheet1", "D", "D", 20)
	f.SetColWidth("Sheet1", "E", "E", 30)
	f.SetColWidth("Sheet1", "F", "F", 20)

	f.SetActiveSheet(index)

	c.Set("Content-Disposition", "attachment; filename=users-report.xlsx")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := f.SaveAs("users-report.xlsx"); err != nil {
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
	var users []m.Users

	result := s.DB.Debug().Find(&users)
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

	for i, user := range users {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), user.IDUser)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), user.Fullname)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), user.Username)
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

	c.Set("Content-Disposition", "attachment; filename=users-report.pdf")
	c.Set("Content-Type", "application/pdf")

	if err := pdf.WritePdf("users-report.pdf"); err != nil {
		log.Println(err)
		return
	}
	errWrite := pdf.Write(c.Writer)
	if errWrite != nil {
		log.Println(errWrite)
		return
	}
}
