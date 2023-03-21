package models

import "time"

type Status string
type PaidStatus string

const (
	New       Status = "baru"
	OnProcess Status = "proses"
	Finished  Status = "selesai"
	Returned  Status = "diambil"

	Paid   PaidStatus = "dibayar"
	UnPaid PaidStatus = "belum_dibayar"
)

type Transactions struct {
	IDTransaction int        `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_transaction"`
	Id_outlet     int        `gorm:"type:INT(11);column:id_outlet" json:"id_outlet"`
	Placement     Outlets    `gorm:"foreignKey:Id_outlet; references:IDOutlet" json:"transaction_at"`
	Invoices      string     `gorm:"unique;column:kode_invoice;type:VARCHAR(100)" json:"invoice"`
	Id_member     int        `gorm:"type:INT(11);column:id_member" json:"member_id"`
	OrderedBy     Members    `gorm:"foreignKey:Id_member;references:IDMember" json:"ordered_by"`
	CreatedDate   time.Time  `gorm:"type:DATETIME;column:tgl" json:"date"`
	Deadline      time.Time  `gorm:"type:DATETIME;column:batas_waktu" json:"deadline"`
	PaidDate      time.Time  `gorm:"type:DATETIME;column:tgl_bayar" json:"paid_date"`
	Addons        int        `gorm:"type:INT(11);column:biaya_tambahan" json:"biaya_tambahan"`
	Discount      int64      `gorm:"type:DOUBLE;column:diskon" json:"discount"`
	Tax           int        `gorm:"type:INT(11);column:pajak" json:"tax"`
	Status        Status     `gorm:"column:status" json:"status"`
	PaidStatus    PaidStatus `gorm:"column:dibayar" json:"paid_status"`
	Id_user       int        `gorm:"type:INT(11);column:id_user" json:"inputter_id"`
	InputBy       Users      `gorm:"foreignKey:Id_user;references:IDUser" json:"input_by"`
}

func (t *Transactions) TableName() string {
	return "tb_transaksi"
}
