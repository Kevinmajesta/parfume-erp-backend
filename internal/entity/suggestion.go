package entity

import (
	"github.com/google/uuid"
)

type Suggestion struct {
	ID_Suggestion uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id_suggestion"`
	UserID        uuid.UUID `json:"user_id"`
	Type          string    `json:"type"`
	Message       string    `json:"message"`
	Auditable
}

func NewSuggestion(tipe, message, date string) *Suggestion {
	return &Suggestion{
		ID_Suggestion: uuid.New(),
		Type:          tipe,
		Message:       message,
		Auditable:     NewAuditable(),
	}
}
