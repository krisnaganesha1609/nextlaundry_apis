package models

type Products struct {
	IDProduct int64   `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_product"`
	Id_outlet int     `gorm:"type:INT(11);column:id_outlet" json:"user_outlet"`
	Outlet    Outlets `gorm:"foreignKey:Id_outlet; references:IDOutlet" json:"outlet"`
	Types     string  `gorm:"type:ENUM('kiloan', 'selimut', 'bed_cover', 'kaos', 'lain'); column:jenis" json:"type"`
	Price     int64   `gorm:"type:INT(11);column:harga" json:"price"`
}

func (m *Products) TableName() string {
	return "tb_paket"
}
