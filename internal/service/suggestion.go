package service

import (
	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/internal/repository"
	"github.com/Kevinmajesta/webPemancingan/pkg/token"
)

type SuggestionService interface {
	CreateSuggestion(suggestion *entity.Suggestion) error
}

type suggestionService struct {
	suggestionRepo repository.SuggestionRepository
	tokenUseCase   token.TokenUseCase
	userRepository repository.UserRepository
}

func NewSuggestionService(suggestionRepo repository.SuggestionRepository, tokenUseCase token.TokenUseCase,
	userRepository repository.UserRepository) SuggestionService {
	return &suggestionService{
		suggestionRepo: suggestionRepo,
		tokenUseCase:   tokenUseCase,
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
