package repository

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type SuggestionRepository interface {
	CreateSuggestion(suggestion *entity.Suggestion) error
}

type suggestionRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewSuggestionRepository(db *gorm.DB, cacheable cache.Cacheable) *suggestionRepository {
	return &suggestionRepository{db: db, cacheable: cacheable}
}

func (r *suggestionRepository) CreateSuggestion(suggestion *entity.Suggestion) error {
	return r.db.Create(suggestion).Error
}
