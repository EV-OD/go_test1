package user

import (
	"context"
	"fmt"
	"myapp/config"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, UID uint) (*UserDTO, error)
	LoadBalanceForUser(ctx context.Context, UID uint, newAmount float64) error
	SendMoney(ctx context.Context, senderId uint, accountId int, amount float64) error
}

type userService struct {
	repo UserRepository
}

func CUUID() int {
	return int(time.Now().Unix())
}

func (u *EUser) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *EUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func GenerateJWT(userID uint) (string, error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	claims := &CustomClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func NewUserService(repo UserRepository) Service {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &EUser{
		UserDTO: UserDTO{
			Name:      req.Name,
			Email:     req.Email,
			Password:  string(hashedPassword),
			AccountID: CUUID(),
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	token, err := GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil

}

func (s *userService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, UID uint) (*UserDTO, error) {
	user, err := s.repo.GetByID(ctx, UID)
	if err != nil {
		return nil, err
	}
	return &user.UserDTO, nil
}

func (s *userService) LoadBalanceForUser(ctx context.Context, UID uint, newAmount float64) error {
	return s.repo.LoadBalanceForUser(ctx, UID, newAmount)
}

func (s *userService) SendMoney(ctx context.Context, senderId uint, accountId int, amount float64) error {
	return s.repo.SendMoney(ctx, senderId, accountId, amount)
}
