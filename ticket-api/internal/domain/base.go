package domain

import (
	"time"

	"ticket-api/internal/types"

	"gorm.io/gorm"
)

type Base struct {
	ID types.BinaryUUID `json:"id"`
}

type BaseDate struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
