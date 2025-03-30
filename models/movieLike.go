package models

import "time"

type MovieLike struct {
	ID       int       `gorm:"primaryKey"`
	MovieID  int       `gorm:"not null;index;foreignKey:ID;constraint:OnDelete:CASCADE" json:"movie_id"`
	UserID   int       `gorm:"not null;index" json:"user_id"`
	CreateAt time.Time `gorm:"autoCreateTime" json:"create_at"`
}
