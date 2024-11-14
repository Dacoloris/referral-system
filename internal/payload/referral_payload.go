package payload

import "time"

type ReferralRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required"`
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	ReferrerID *uint  `json:"referrer_id"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateRequest struct {
	UserID     uint      `json:"user_id" binding:"required"`
	Email      string    `json:"email" binding:"required"`
	ExpiryDate time.Time `json:"expiry_date"`
}
