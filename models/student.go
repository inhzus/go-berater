package models

type Student struct {
	Openid string `gorm:"type:varchar(30);primary_key"`
	Phone string `gorm:"type:varchar(20);not null"`
	IdCard string `gorm:"column:id_card;type:varchar(18);not null"`
}
