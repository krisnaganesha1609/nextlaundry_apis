package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	Admin   Role = "admin"
	Cashier Role = "kasir"
	Owner   Role = "owner"
)

type Users struct {
	IDUser    int     `gorm:"type:INT(11);primaryKey;column:id;autoIncrement" json:"id_user"`
	Fullname  string  `gorm:"type:VARCHAR(100);column:nama" json:"fullname"`
	Username  string  `gorm:"type:VARCHAR(40);column:username" json:"username"`
	Password  string  `gorm:"type:VARCHAR(100);column:password" json:"password"`
	Role      Role    `gorm:"type:ENUM('admin', 'kasir', 'owner');column:role" json:"role"`
	Id_outlet int     `gorm:"type:INT(11);column:id_outlet" json:"user_outlet"`
	Placement Outlets `gorm:"foreignKey:Id_outlet; references:IDOutlet" json:"placement"`
}

func (u *Users) TableName() string {
	return "tb_user"
}

func (u *Users) HashingPassword(password string) (hashed string, err error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedByte)

	return hashedPassword, nil
}

func (u *Users) CheckPasswordHash(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) PrepareGive() {
	u.Password = ""
}
