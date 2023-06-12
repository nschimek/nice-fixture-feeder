package model

const (
	StatusTypeInPlay = "IP"
	StatusTypeFinished = "FI"
	StatusTypeScheduled = "SC"
)

type FixtureStatus struct {
	Id string `gorm:"primaryKey"`
	Name, Type, Description string
}