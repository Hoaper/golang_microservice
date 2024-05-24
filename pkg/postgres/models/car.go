package models

import (
	"time"
)

type Car struct {
	ID        uint   `gorm:"primary_key"`
	RegNum    string `gorm:"type:varchar(20);unique_index"`
	Mark      string `gorm:"type:varchar(100)"`
	Model     string `gorm:"type:varchar(100)"`
	Year      int
	OwnerID   uint
	Owner     People
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
