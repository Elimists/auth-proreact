package models

import (
	"time"
)

// CustomModel is used instead of gorm.Model to avoid the DeletedAt fields.
type CustomModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	CustomModel
	Email            string           `json:"email" gorm:"unique;primary_key"` // The email address of the user. This is the primary key.
	Password         []byte           `json:"-"`
	Privilege        int8             `json:"privilege"` // 1: Admin, 2: Manager, 3: Coordinator, 4: Moderator, 9: General user
	Verified         bool             `json:"-"`
	UserVerification UserVerification `json:"userVerification" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}

type UserVerification struct {
	CustomModel
	UserID             uint   `json:"-"` // The user ID of the user this verification belongs to.
	VerificationCode   string `json:"-"`
	VerificationExpiry uint   `json:"-"`
}
