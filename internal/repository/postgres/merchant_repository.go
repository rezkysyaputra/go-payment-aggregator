package postgres

import (
	"context"
	"go-payment-aggregator/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"size:255;not null"`
	Email       string    `gorm:"size:255;unique;not null"`
	ApiKey      string    `gorm:"size:255;unique"`
	CallbackURL string    `gorm:"size:255"`
	Status      string    `gorm:"size:50;not null;default:'ACTIVE'"`
	Balance     int64     `gorm:"default:0;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (MerchantModel) TableName() string {
	return "merchants"
}

// toMerchantModel converts domain.Merchant to MerchantModel
func toMerchantModel(d *domain.Merchant) *MerchantModel {
	return &MerchantModel{
		ID:          d.ID,
		Name:        d.Name,
		Email:       d.Email,
		ApiKey:      d.APIKeyHash,
		CallbackURL: d.CallbackURL,
		Status:      string(d.Status),
		Balance:     d.Balance,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// toDomain converts MerchantModel to domain.Merchant
func (m *MerchantModel) toDomain() *domain.Merchant {
	return &domain.Merchant{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		APIKeyHash:  m.ApiKey,
		CallbackURL: m.CallbackURL,
		Status:      domain.MerchantStatus(m.Status),
		Balance:     m.Balance,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) domain.MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

// Create inserts a new merchant into the database
func (r *merchantRepository) Create(ctx context.Context, m *domain.Merchant) error {
	model := toMerchantModel(m)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

// Update modifies an existing merchant in the database
func (r *merchantRepository) Update(ctx context.Context, m *domain.Merchant) error {
	model := toMerchantModel(m)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

// FindByApiKey retrieves a merchant by its API key
func (r *merchantRepository) FindByApiKey(ctx context.Context, apiKey string) (*domain.Merchant, error) {
	var model MerchantModel
	if err := r.db.WithContext(ctx).First(&model, "api_key = ?", apiKey).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

// FindByID retrieves a merchant by its ID
func (r *merchantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	var model MerchantModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

// RegenerateApiKey updates the API key for a merchant
func (r *merchantRepository) RegenerateApiKey(ctx context.Context, id uuid.UUID, newApiKey string) error {
	if err := r.db.WithContext(ctx).Model(MerchantModel{}).Where("id = ?", id).Update("api_key = ?", newApiKey).Error; err != nil {
		return err
	}
	return nil
}
