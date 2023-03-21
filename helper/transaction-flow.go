package helper

import (
	"log"
	"net/http"
	m "nextlaundry_apis/models"
	s "nextlaundry_apis/models/setup"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func CreateViewToExport(c *gin.Context) {
	//TODO: CREATE VIEW TRANSACTION IN HERE!

	s.DB.Exec("CREATE OR REPLACE VIEW `transaction_report` AS SELECT tb_transaksi.kode_invoice AS invoice, tb_member.nama AS member_name, tb_transaksi.tgl AS date, tb_transaksi.tgl_bayar AS paid_date, tb_transaksi.biaya_tambahan AS add_charge, tb_transaksi.diskon AS discount, tb_transaksi.pajak AS tax, tb_outlet.nama AS transacted_at, tb_transaksi.dibayar AS paid_status, tb_user.nama AS cashier, tb_paket.jenis AS laundry_type, tb_paket.harga AS unit_price, tb_detail_transaksi.qty AS qty FROM (((((tb_transaksi INNER JOIN tb_member ON tb_transaksi.id_member = tb_member.id) INNER JOIN tb_user ON tb_transaksi.id_user = tb_user.id) INNER JOIN tb_detail_transaksi ON tb_transaksi.id = tb_detail_transaksi.id_transaksi) INNER JOIN tb_paket ON tb_detail_transaksi.id_paket = tb_paket.id) INNER JOIN tb_outlet ON tb_transaksi.id_outlet = tb_outlet.id);  ")

	sqlQuery := "SELECT * FROM transaction_report"

	viewName := "transaction_report"
	// Create a temporary table to hold the result of the SQL query
	tmpTableName := "tmp_" + viewName
	exists := s.DB.Exec("CREATE OR REPLACE TEMPORARY TABLE " + tmpTableName + " AS " + sqlQuery)

	if err := exists.Error; err != nil {
		s.DB.Exec("DROP IF EXISTS " + tmpTableName)
	}

	// Map the temporary table to a GORM model
	s.DB.Table(tmpTableName).Find(m.TransactionReport{})
}

func ExportViewToExcel(c *gin.Context) {
	CreateViewToExport(c)
	var myViews []m.TransactionReport
	viewName := "transaction_report"
	tmpTableName := "tmp_" + viewName
	results := s.DB.Table(tmpTableName).Debug().Find(&myViews)

	if results.Error != nil {
		log.Println(results.Error)
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
	f.SetCellValue("Sheet1", "B1", "Invoice")
	f.SetCellValue("Sheet1", "C1", "Member")
	f.SetCellValue("Sheet1", "D1", "Date")
	f.SetCellValue("Sheet1", "E1", "Paid Date")
	f.SetCellValue("Sheet1", "F1", "Add Charge")
	f.SetCellValue("Sheet1", "G1", "Discount")
	f.SetCellValue("Sheet1", "H1", "Tax")
	f.SetCellValue("Sheet1", "I1", "Paid Status")
	f.SetCellValue("Sheet1", "J1", "Cashier")
	f.SetCellValue("Sheet1", "K1", "Laundry Type")
	f.SetCellValue("Sheet1", "L1", "Unit Price")
	f.SetCellValue("Sheet1", "M1", "Quantity")

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

	f.SetCellStyle("Sheet1", "A1", "M1", style)

	for i, v := range myViews {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), v.Invoice)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), v.Member)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), v.Date)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), v.PaidDate)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), v.AddCharge)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), v.Discount)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), v.Tax)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), v.PaidStatus)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i+2), v.Cashier)
		f.SetCellValue("Sheet1", "K"+strconv.Itoa(i+2), v.LaundryType)
		f.SetCellValue("Sheet1", "L"+strconv.Itoa(i+2), v.UnitPrice)
		f.SetCellValue("Sheet1", "M"+strconv.Itoa(i+2), v.Qty)
	}

	f.SetColWidth("Sheet1", "A", "A", 5)
	f.SetColWidth("Sheet1", "B", "B", 30)
	f.SetColWidth("Sheet1", "C", "C", 30)
	f.SetColWidth("Sheet1", "D", "D", 30)
	f.SetColWidth("Sheet1", "E", "E", 30)
	f.SetColWidth("Sheet1", "F", "F", 30)
	f.SetColWidth("Sheet1", "G", "G", 30)
	f.SetColWidth("Sheet1", "H", "H", 30)
	f.SetColWidth("Sheet1", "I", "I", 30)
	f.SetColWidth("Sheet1", "J", "J", 30)
	f.SetColWidth("Sheet1", "K", "K", 30)
	f.SetColWidth("Sheet1", "L", "L", 30)
	f.SetColWidth("Sheet1", "M", "M", 5)

	f.SetActiveSheet(index)

	c.Set("Content-Disposition", "attachment; filename=transaction-report.xlsx")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := f.SaveAs("transaction-report.xlsx"); err != nil {
		log.Println(err)
		return
	}
	errWrite := f.Write(c.Writer)
	if errWrite != nil {
		log.Println(errWrite)
		return
	}

}

// func ExportViewToPDF(c *gin.Context) {
// 	CreateViewToExport(c)
// 	var myViews []m.LogTransaction
// 	viewName := "transaction_report"
// 	tmpTableName := "tmp_" + viewName
// 	result := s.DB.Table(tmpTableName).Debug().Find(&myViews)
// 	if result.Error != nil {
// 		log.Println(result.Error)
// 	}

// 	f := excelize.NewFile()

// 	index, err := f.NewSheet("Sheet1")

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	f.SetCellValue("Sheet1", "A1", "No")
// 	f.SetCellValue("Sheet1", "B1", "Invoice")
// 	f.SetCellValue("Sheet1", "C1", "Member")
// 	f.SetCellValue("Sheet1", "D1", "Date")
// 	f.SetCellValue("Sheet1", "E1", "Paid Date")
// 	f.SetCellValue("Sheet1", "F1", "Add Charge")
// 	f.SetCellValue("Sheet1", "G1", "Discount")
// 	f.SetCellValue("Sheet1", "H1", "Tax")
// 	f.SetCellValue("Sheet1", "I1", "Paid Status")
// 	f.SetCellValue("Sheet1", "J1", "Cashier")
// 	f.SetCellValue("Sheet1", "K1", "Laundry Type")
// 	f.SetCellValue("Sheet1", "L1", "Unit Price")
// 	f.SetCellValue("Sheet1", "M1", "Quantity")

// 	style, err := f.NewStyle(&excelize.Style{
// 		Font: &excelize.Font{
// 			Bold: true,
// 		},
// 	})

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	f.SetCellStyle("Sheet1", "A1", "M1", style)

// 	for i, v := range myViews {
// 		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
// 		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), v.Invoice)
// 		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), v.Member)
// 		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), v.Date)
// 		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), v.PaidDate)
// 		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), v.AddCharge)
// 		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), v.Discount)
// 		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), v.Tax)
// 		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), v.PaidStatus)
// 		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i+2), v.Cashier)
// 		f.SetCellValue("Sheet1", "K"+strconv.Itoa(i+2), v.LaundryType)
// 		f.SetCellValue("Sheet1", "L"+strconv.Itoa(i+2), v.UnitPrice)
// 		f.SetCellValue("Sheet1", "M"+strconv.Itoa(i+2), v.Qty)
// 	}

// 	f.SetColWidth("Sheet1", "A", "A", 5)
// 	f.SetColWidth("Sheet1", "B", "B", 30)
// 	f.SetColWidth("Sheet1", "C", "C", 30)
// 	f.SetColWidth("Sheet1", "D", "D", 30)
// 	f.SetColWidth("Sheet1", "E", "E", 30)
// 	f.SetColWidth("Sheet1", "F", "F", 30)
// 	f.SetColWidth("Sheet1", "G", "G", 30)
// 	f.SetColWidth("Sheet1", "H", "H", 30)
// 	f.SetColWidth("Sheet1", "I", "I", 30)
// 	f.SetColWidth("Sheet1", "J", "J", 30)
// 	f.SetColWidth("Sheet1", "K", "K", 30)
// 	f.SetColWidth("Sheet1", "L", "L", 30)
// 	f.SetColWidth("Sheet1", "M", "M", 5)

// 	f.SetActiveSheet(index)

// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

// 	errFont := pdf.AddTTFFont("righteous", "C:/Backup Data/Project/Rich Go/nextlaundry_apis/asset/fonts/Righteous-Regular.ttf")
// 	if errFont != nil {
// 		log.Println("failed to add font")
// 	}
// 	errFont = pdf.SetFont("righteous", "", 14)
// 	if errFont != nil {
// 		log.Println("failed to set font")
// 	}

// 	pdf.AddPage()

// 	r, err := f.GetRows("Sheet1")
// 	for row, rowCells := range r {
// 		for _, cell := range rowCells {

// 			err = pdf.Cell(nil, cell)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			pdf.SetX(pdf.GetX() + 100)
// 		}

// 		pdf.Br(30)
// 		pdf.SetX(20)

// 		if row%20 == 19 {
// 			pdf.AddPage()
// 			pdf.SetX(20)
// 		}

// 	}

// 	c.Set("Content-Disposition", "attachment; filename=transaction-report.pdf")
// 	c.Set("Content-Type", "application/pdf")

// 	if err := pdf.WritePdf("transaction-report.pdf"); err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	errWrite := pdf.Write(c.Writer)
// 	if errWrite != nil {
// 		log.Println(errWrite)
// 		return
// 	}
// }

func GetDashboardNeeds() {
	//TODO: CREATE DASHBOARD FUNCTION BERDASARKAN TRANSACTION IN HERE
}
