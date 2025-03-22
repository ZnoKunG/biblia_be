package model

type Book struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	NumPages int     `json:"nopages"`
	Language string  `json:"lang"`
	Price    float32 `json:"price"`
	AuthorID uint    `json:"authorID"`
}
