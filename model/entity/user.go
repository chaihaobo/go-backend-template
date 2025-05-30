package entity

import (
	"github.com/chaihaobo/be-template/utils/jwt"
)

type User struct {
	BaseEntity
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}

func (u User) ToJWTClaims() *jwt.UserForToken {
	return &jwt.UserForToken{
		ID: u.ID,
	}
}
