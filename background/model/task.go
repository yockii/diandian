package model

type Task struct {
	Base
	Name        string `json:"name" gorm:"size:200"`
	Description string `json:"description"`
}
