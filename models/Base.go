package models

type Model struct {
	ID int `gorm:"primary_key" form:"id" json:"id"`
}
