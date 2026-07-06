package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *EUser) error
	GetByEmail(ctx context.Context, email string) (*EUser, error)
	GetByID(ctx context.Context, ID uint) (*EUser, error)
	LoadBalanceForUser(ctx context.Context, UID uint, newAmount float64) error
	SendMoney(ctx context.Context, senderId uint, accountId int, amount float64) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) UserRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, user *EUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*EUser, error) {
	var user EUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, ID uint) (*EUser, error) {
	var user EUser
	if err := r.db.WithContext(ctx).First(&user, ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) LoadBalanceForUser(ctx context.Context, UID uint, newAmount float64) error {

	err := r.db.WithContext(ctx).Model(&EUser{}).Where("id = ?", UID).Update("balance", gorm.Expr("balance + ?", newAmount))

	if err != nil {
		return err.Error
	}
	return nil
}

func (r *postgresRepository) SendMoney(ctx context.Context, senderId uint, accountId int, amount float64) error {
	var sender EUser
	var receiver EUser

	if err := r.db.WithContext(ctx).First(&sender, senderId).Error; err != nil {
		return err
	}

	if sender.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	if err := r.db.WithContext(ctx).Where("account_id = ?", accountId).First(&receiver).Error; err != nil {
		return err
	}

	tx := r.db.WithContext(ctx).Begin()
	txErr := tx.Model(&sender).Where("id = ?", senderId).Update("balance", gorm.Expr("balance - ?", amount)).Error
	if txErr != nil {
		tx.Rollback()
		return txErr
	}

	txErr = tx.Model(&receiver).Where("account_id = ?", accountId).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if txErr != nil {
		tx.Rollback()
		return txErr
	}

	tx.Commit()
	return nil

}
