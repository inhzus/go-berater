package models

import "github.com/jinzhu/gorm"

type Student struct {
	gorm.Model
	Openid string `gorm:"type:varchar(30);not null;index"`
	Phone string `gorm:"type:varchar(20);not null"`
	IdCard string `gorm:"column:id_card;type:varchar(18);not null"`
}
