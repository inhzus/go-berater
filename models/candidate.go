package models

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
	Openid   string `gorm:"type:varchar(30);primary_key"`
	Phone    string `gorm:"type:varchar(20);not null"`
	Name     string `gorm:"type:varchar(16);not null"`
	Province string `gorm:"type:varchar(2);not null"`
	City     string `gorm:"type:varchar(10);not null"`
	Score    string `gorm:"type:float(4,1);not null"`
	Subject  string `gorm:"type:varchar(10);not null"`
}

func ExistCandidateById(openid string) bool {
	var candidate Candidate
	db.Select("openid").Where("openid = ?", openid).First(&candidate)
	return candidate.Openid == openid
}

func AddCandidate(candidate *Candidate) error {
	return db.Create(candidate).Error
}
