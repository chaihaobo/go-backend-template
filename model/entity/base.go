package entity

import "time"

type BaseEntity struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
