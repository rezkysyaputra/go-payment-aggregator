package postgres

import (
	"context"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	MerchantID    uuid.UUID `gorm:"type:uuid;not null"`
	OrderID       string    `gorm:"size:255;not null;unique"`
	Provider      string    `gorm:"size:255"`
	PaymentMethod string    `gorm:"size:255"`
	Amount        int64     `gorm:"default:0;not null"`
	Currency      string    `gorm:"size:10;not null;default:'IDR'"`
	Status        string    `gorm:"size:50;not null;default:'PENDING'"`
	ExternalRef   string    `gorm:"size:255;not null"`
	RedirectURL   string    `gorm:"size:255"`
	RawResponse   []byte    `gorm:"type:jsonb"`
	ExpiredAt     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (TransactionModel) TableName() string {
	return "transactions"
}

// toTransactionModel converts domain.Transaction to TransactionModel
func toTransactionModel(tx *domain.Transaction) *TransactionModel {
	return &TransactionModel{
		ID:            tx.ID,
		MerchantID:    tx.MerchantID,
		OrderID:       tx.OrderID,
		Provider:      tx.Provider,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Status:        string(tx.Status),
		ExternalRef:   tx.ExternalID,
		PaymentMethod: tx.PaymentMethod,
		RedirectURL:   tx.PaymentURL,
		RawResponse:   pkg.JsonToByte(tx.RawResponse),
		ExpiredAt:     tx.ExpiredAt,
		CreatedAt:     tx.CreatedAt,
		UpdatedAt:     tx.UpdatedAt,
	}
}

// toDomain converts TransactionModel to domain.Transaction
func (t *TransactionModel) toDomain() *domain.Transaction {
	return &domain.Transaction{
		ID:            t.ID,
		MerchantID:    t.MerchantID,
		OrderID:       t.OrderID,
		Provider:      t.Provider,
		Amount:        t.Amount,
		Currency:      t.Currency,
		Status:        domain.TransactionStatus(t.Status),
		ExternalID:    t.ExternalRef,
		PaymentMethod: t.PaymentMethod,
		PaymentURL:    t.RedirectURL,
		RawResponse:   string(t.RawResponse),
		ExpiredAt:     t.ExpiredAt,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (t *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	model := toTransactionModel(tx)
	model.RawResponse = nil
	if err := t.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (t *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	model := toTransactionModel(tx)

	updateData := map[string]interface{}{
		"redirect_url": model.RedirectURL,
		"external_ref": model.ExternalRef,
		"raw_response": model.RawResponse,
		"expired_at":   model.ExpiredAt,
	}

	if err := t.db.WithContext(ctx).Model(&TransactionModel{}).Where("id = ?", model.ID).Updates(updateData).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (t *transactionRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var model TransactionModel
	if err := t.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (t *transactionRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Transaction, error) {
	var model TransactionModel
	if err := t.db.WithContext(ctx).Where("order_id = ?", orderID).First(&model).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}
