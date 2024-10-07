package binder

type CreateSuggestion struct {
	UserId  string `json:"user_id" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Message string `json:"message" validate:"required"`
}
