package models

import (
	"time"
)

// People model
type People struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"type:varchar(100)"`
	Surname   string `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
