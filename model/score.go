package model

type Score struct {
	Id int `gorm:"primaryKey"`
	Name, Description string
	Weight float64
}