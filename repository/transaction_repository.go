package repository

import (
	"context"
	"errors"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TransactionRepository interface {
		GetTransactionByReference(ctx context.Context, tx *gorm.DB, reference string) (entity.Transaction, error)
		UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) error
		SoftDeleteTransaction(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	}

	transactionRepository struct {
		db *gorm.DB
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetTransactionByReference(ctx context.Context, tx *gorm.DB, reference string) (entity.Transaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.Transaction
	if err := tx.WithContext(ctx).Where("reference = ?", reference).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Transaction{}, dto.ErrTransactionNotFound
		}
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&transaction).Error; err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) SoftDeleteTransaction(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Transaction{}).Error
}
