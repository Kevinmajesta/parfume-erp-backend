package repository

import (
	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/pkg/cache"
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
