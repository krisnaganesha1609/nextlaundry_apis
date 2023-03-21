package models

type Outlets struct {
	IDOutlet   int64  `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_outlet"`
	OutletName string `gorm:"type:VARCHAR(100);column:nama" json:"nama_outlet"`
	Address    string `gorm:"type:text;column:alamat" json:"alamat"`
	Phone      string `gorm:"type:VARCHAR(15);column:telepon" json:"telepon"`
	TotalEmp   int64  `gorm:"type:INT(10);column:jumlah_user" json:"total_emp"`
}

func (m *Outlets) TableName() string {
	return "tb_outlet"
}
