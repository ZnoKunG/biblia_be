package model

type User struct {
	ID             uint     `gorm:"primaryKey" json:"id"`
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	FavoriteGenres []string `json:"favorite_genres" gorm:"serializer:json"`
	Records        []Record `json:"records" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreateUser struct {
	Username       string   `json:"username" binding:"required"`
	Password       string   `json:"password" binding:"required"`
	FavoriteGenres []string `json:"favorite_genres"`
}

type UpdateUser struct {
	Username       string   `json:"username" binding:"required"`
	Password       string   `json:"password" binding:"required"`
	FavoriteGenres []string `json:"favorite_genres"`
}
