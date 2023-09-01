package model

type Score struct {
	Id int `gorm:"primaryKey"`
	Name, Description string
	Weight int
}

type FixtureScore struct {
	FixtureId int `gorm:"primaryKey"`
	ScoreId int `gorm:"primaryKey"`
	Value int
	Audit
}