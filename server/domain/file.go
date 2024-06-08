package domain

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	gorm.Model
	Name      string    `gorm:"not null"`
	Url       string    `gorm:"not null"`
	UrlExpiry time.Time `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	//RecipeID uint `gorm:""`
}

type FileDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (f *File) ToDto() Dto {
	return FileDto{
		ID:   f.ID,
		Name: f.Name,
		Url:  f.Url,
	}
}
