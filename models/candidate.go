package models

import "github.com/jinzhu/gorm"

/*
   openid = engine.Column(engine.String(30), primary_key=True)
   phone = engine.Column(engine.String(20), nullable=False)
   name = engine.Column(engine.String(16), nullable=False)
   province = engine.Column(engine.String(2), nullable=False)
   city = engine.Column(engine.String(10), nullable=False)
   score = engine.Column(engine.Float(precision=1), nullable=False)
   subject = engine.Column(engine.String(10), nullable=False)
*/
type Candidate struct {
	gorm.Model
	Openid   string `gorm:"type:varchar(30);not null;index"`
	Phone    string `gorm:"type:varchar(20);not null"`
	Name     string `gorm:"type:varchar(16);not null"`
	Province string `gorm:"type:varchar(2);not null"`
	City     string `gorm:"type:varchar(10);not null"`
	Score    string `gorm:"type:float(4,1);not null"`
	Subject  string `gorm:"type:varchar(10);not null"`
}

func ExistCandidateById(openid string) bool {
	var count int
	db.Model(&Candidate{}).Where("openid = ?", openid).Count(&count)
	return count != 0
}

func GetCandidateById(openid string) *Candidate {
	var candidate []Candidate
	db.Where("openid = ?", openid).Find(&candidate)
	if len(candidate) != 0 {
		return &candidate[0]
	} else {
		return nil
	}
}

func AddCandidate(candidate *Candidate) error {
	return db.Create(candidate).Error
}

func UpdateCandidate(openid string, candidate *Candidate) error {
	return db.Model(&Candidate{}).Where("openid = ?", openid).Update(candidate).Error
}

func RemoveCandidateById(openid string) bool {
	err := db.Where("openid = ?", openid).Delete(&Candidate{}).Error
	return err == nil
}
