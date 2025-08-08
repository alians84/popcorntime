package models

import "time"

type FilesS3 struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	FileURL   string    `json:"file_url" gorm:"size:255;unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}
