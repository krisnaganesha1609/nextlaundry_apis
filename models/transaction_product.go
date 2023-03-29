package models

type TransactionProduct struct {
	ID              int          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Id_transaksi    int          `gorm:"type:INT(11);column:id_transaksi" json:"transaction_id"`
	TransactionInfo Transactions `gorm:"foreignKey:Id_transaksi;references:IDTransaction" json:"-"`
	Id_product      int          `gorm:"type:INT(11);column:id_product" json:"id_product"`
	ProductInfo     Products     `gorm:"foreignKey:Id_product;references:IDProduct" json:"-"`
	Qty             int          `gorm:"type:INT(11);column:qty" json:"qty"`
}

func (t *TransactionProduct) TableName() string {
	return "tb_transaksi_produk"
}
