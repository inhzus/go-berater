package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

func Setup() {
	var err error
	db, err = gorm.Open("mysql", "root:root@/berater?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatalf("Models.setup err: %v", err)
	}
	//db.CreateTable(&Candidate{})
	//db.CreateTable(&Student{})
	//db.CreateTable(&Freshman{})
	db.AutoMigrate(&Candidate{})
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Freshman{})
}

func Teardown() {
	_ = db.Close()
}
