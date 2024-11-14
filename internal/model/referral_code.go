package model

import "time"

type ReferralCode struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Code       string    `json:"code" gorm:"unique"`
	UserID     uint      `json:"user_id" gorm:"unique"`
	ExpiryDate time.Time `json:"expiry_date"`
	CreatedAt  time.Time `json:"created_at"`
}
