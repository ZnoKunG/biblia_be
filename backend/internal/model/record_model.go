package model

import (
	"time"
)

type Record struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID      uint      `json:"userID" gorm:"uniqueIndex:idx_user_isbn;foreignKey:UserID;references:ID"`
	ISBN        string    `json:"isbn" gorm:"type:varchar(20);uniqueIndex:idx_user_isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Cover       string    `json:"cover"`
	Genre       string    `json:"genre"`
	Status      string    `json:"status"`
	CurrentPage int32     `json:"currentPage"`
	TotalPages  int32     `json:"totalPages"`
	DateAdded   time.Time `json:"dateAdded"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateRecord struct {
	UserID      uint      `json:"userID"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Cover       string    `json:"cover"`
	Genre       string    `json:"genre"`
	Status      string    `json:"status"`
	CurrentPage int32     `json:"currentPage"`
	TotalPages  int32     `json:"totalPages"`
	DateAdded   time.Time `json:"dateAdded"`
}

type UpdateRecord struct {
	Status      string `json:"status"`
	CurrentPage int32  `json:"currentPage"`
}
