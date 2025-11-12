package transaction

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *Transaction) error
	UpdateStatus(id uuid.UUID, status string) error
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

// Create implements TransactionRepository.
func (r *TransactionRepositoryImpl) Create(tx *Transaction) error {
	return r.db.Create(tx).Error
}

// UpdateStatus implements TransactionRepository.
func (r *TransactionRepositoryImpl) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&Transaction{}).
		Where("id = ?", id).
		Update("status", status).Error
}
