package models

// type TransactionDetails struct {
// 	IDDetail        int          `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_detail"`
// 	Id_transaksi    int          `gorm:"type:INT(11);column:id_transaksi" json:"transaction_id"`
// 	TransactionInfo Transactions `gorm:"foreignKey:Id_transaksi;references:IDTransaction" json:"transaction_info"`
// 	Id_paket        int          `gorm:"type:INT(11);column:id_paket" json:"id_paket"`
// 	Packages        Products     `gorm:"foreignKey:Id_paket;references:IDProduct" json:"package_details"`
// 	Quantity        int64        `gorm:"type:DOUBLE;column:qty" json:"qty"`
// 	Desc            string       `gorm:"type:TEXT;column:keterangan" json:"description"`
// }
type TransactionDetails struct {
	IDDetail        int          `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_detail"`
	Id_transaksi    int          `gorm:"type:INT(11);column:id_transaksi" json:"transaction_id"`
	TransactionInfo Transactions `gorm:"foreignKey:Id_transaksi;references:IDTransaction" json:"-"`
}

func (td *TransactionDetails) TableName() string {
	return "tb_detail_transaksi"
}
