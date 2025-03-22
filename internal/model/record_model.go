package model

import (
	"time"
)

type Record struct {
	BookID        uint       `json:"bookID" gorm:"primaryKey;foreignKey:BookID;references:ID"`
	UserID        uint       `json:"userID" gorm:"primaryKey;foreignKey:UserID;references:ID"`
	Status        string     `json:"status"`
	Curr_progress int32      `json:"curr_progress"`
	Curr_chapter  int32      `json:"curr_chapter"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedDate   *time.Time `json:"started_date"`
	UpdateDate    *time.Time `json:"update_date"`
	StopDate      *time.Time `json:"stop_date"`
	FinishDate    *time.Time `json:"finish_date"`

	Book Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateRecord struct {
	BookID        uint   `json:"bookID"`
	UserID        uint   `json:"userID"`
	Status        string `json:"status"`
	Curr_progress int32  `json:"curr_progress"`
	Curr_chapter  int32  `json:"curr_chapter"`
}

type UpdateRecord struct {
	Status        string `json:"status"`
	Curr_progress int32  `json:"curr_progress"`
	Curr_chapter  int32  `json:"curr_chapter"`
}
