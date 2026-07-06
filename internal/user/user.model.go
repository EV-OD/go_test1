package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID        uint
	Name      string  `json:"name" gorm:"type:varchar(100);not null"`
	Email     string  `json:"email" gorm:"type:varchar(100);unique;not null"`
	AccountID int     `json:"account_id" gorm:"unique;not null"`
	Balance   float64 `json:"balance" gorm:"type:decimal(10,2);default:0.00"`
	Password  string  `json:"-" gorm:"type:varchar(255);not null"`
}

type EUser struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	UserDTO
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
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
	Amount float64 `json:"amount" binding:"required"`
}

type SendMoneyRequest struct {
	ReceiverAccountID int     `json:"receiver_account_id" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
}
