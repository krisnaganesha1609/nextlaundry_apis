package models

type Members struct {
	IDMember      int64  `gorm:"type:INT(11);primaryKey;autoIncrement;column:id" json:"id_member"`
	MemberName    string `gorm:"column:nama;type:VARCHAR(100)" json:"member_name"`
	MemberAddress string `gorm:"column:alamat;type:text" json:"member_address"`
	Gender        string `gorm:"column:jenis_kelamin;type:ENUM(L, P)" json:"gender"`
	MemberPhone   string `gorm:"column:telepon;type:VARCHAR(15)" json:"member_phone"`
}

func (m *Members) TableName() string {
	return "tb_member"
}
