package models

import "github.com/jinzhu/gorm"

type Freshman struct {
	gorm.Model
	Stuid string `gorm:"type:varchar(12);not null;index"`
	Name string `gorm:"type:varchar(16)"`
	Origin string `gorm:"type:varchar(10);index"`
	Gender string `gorm:"type:varchar(2)"`
	Department string `gorm:"type:varchar(32);index"`
	IdCard string `gorm:"type:varchar(20);index;not null"`
	AdmissionId string `gorm:"column:admission_id;type:varchar(20);index"`
}
