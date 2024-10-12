package service

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type SuggestionService interface {
	CreateSuggestion(suggestion *entity.Suggestion) error
}

type suggestionService struct {
	suggestionRepo repository.SuggestionRepository
	userRepository repository.UserRepository
}

func NewSuggestionService(suggestionRepo repository.SuggestionRepository, 
	userRepository repository.UserRepository) SuggestionService {
	return &suggestionService{
		suggestionRepo: suggestionRepo,
		userRepository: userRepository,
	}
}

func (s *suggestionService) CreateSuggestion(suggestion *entity.Suggestion) error {
	userIds, err := s.userRepository.GetAllUserIds()
	if err != nil {
		return err
	}

	for _, userID := range userIds {
		// Buat salinan notifikasi untuk setiap pengguna
		userSuggestion := *suggestion
		userSuggestion.UserID = userID
		if err := s.suggestionRepo.CreateSuggestion(&userSuggestion); err != nil {
			return err
		}
	}

	return nil
}
