package models

type TransactionDetailsRequest struct {
	Id_transaksi int `gorm:"type:INT(11);column:id_transaksi" json:"id_transaksi"`
	Id_product   int `json:"id_product"`
	Qty          int `json:"qty"`
}
