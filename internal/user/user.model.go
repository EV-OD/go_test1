package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID        uint    `example:"1"`
	Name      string  `json:"name" gorm:"type:varchar(100);not null" example:"John Doe"`
	Email     string  `json:"email" gorm:"type:varchar(100);unique;not null" example:"john@example.com"`
	AccountID int     `json:"account_id" gorm:"unique;not null" example:"1001"`
	Balance   float64 `json:"balance" gorm:"type:decimal(10,2);default:0.00" example:"250.50"`
	Password  string  `json:"-" gorm:"type:varchar(255);not null"`
	Role      string  `json:"role" gorm:"type:varchar(20);check:role IN ('viewer', 'editor');default:'viewer';not null" example:"viewer"`
}

type EUser struct {
	CreatedAt time.Time      `example:"2025-01-15T10:30:00Z"`
	UpdatedAt time.Time      `example:"2025-06-01T14:20:00Z"`
	DeletedAt gorm.DeletedAt
	UserDTO
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"mypassword"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"mypassword"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  EUser  `json:"user"`
}

type UserJwtParsedData struct {
	Token string `json:"token" binding:"required"`
	User  EUser  `json:"user" binding:"required"`
}

type CustomClaims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}

type LoadBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required" example:"100.00"`
}

type SendMoneyRequest struct {
	ReceiverAccountID int     `json:"receiver_account_id" binding:"required" example:"1002"`
	Amount            float64 `json:"amount" binding:"required" example:"50.00"`
}
