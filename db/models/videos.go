package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// Videos model entity
type Videos struct {
	Id          uuid.UUID      `json:"id" gorm:"primaryKey"`
	YoutubeId   sql.NullString `json:"youtubeId" gorm:"unique"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   time.Time      `json:"deletedAt" gorm:"default:null"`
}
