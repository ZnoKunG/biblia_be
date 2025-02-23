package model

import (
	"time"
)

type Record struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	BookID        uint       `json:"bookID"`
	UserID        uint       `json:"userID"`
	Status        string     `json:"status"`
	Curr_progress int32      `json:"curr_progress"`
	Curr_chapter  int32      `json:"curr_chapter"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedDate   *time.Time `json:"started_date"`
	UpdateDate    *time.Time `json:"update_date"`
	StopDate      *time.Time `json:"stop_date"`
	FinishDate    *time.Time `json:"finish_date"`
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
