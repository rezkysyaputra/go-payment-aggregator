package transaction

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *Transaction) error
	UpdateStatus(id uuid.UUID, status string) error
	FindById(id uuid.UUID) (*Transaction, error)
	FindOrderById(orderID string) (*Transaction, error)
	UpdateStatusAndRaw(id uuid.UUID, status string, rawJSON string) error
}

type TransactionRepositoryImpl struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewTransactionRepository(db *gorm.DB, log *logrus.Logger) TransactionRepository {
	return &TransactionRepositoryImpl{db: db, log: log}
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

// FindById implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindById(id uuid.UUID) (*Transaction, error) {
	var t Transaction

	// find by id
	err := r.db.Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// FindOrderById implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindOrderById(orderID string) (*Transaction, error) {
	var t Transaction

	// find by orderID
	err := r.db.Where("order_id = ?", orderID).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// UpdateStatusAndRaw implements TransactionRepository.
func (r *TransactionRepositoryImpl) UpdateStatusAndRaw(id uuid.UUID, status string, rawJSON string) error {
	// update status and raw response
	return r.db.Model(&Transaction{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       status,
			"raw_response": rawJSON,
			"updated_at":   gorm.Expr("NOW()"),
		}).Error
}
