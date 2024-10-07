package binder

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminCreateRequest struct {
	Fullname     string `json:"fullname" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	Role         string `json:"role" `
	Phone        string `json:"phone" validate:"required"`
	Status       bool   `json:"status" `
	Verification bool   `json:"verification" `
}

type AdminUpdateRequest struct {
	Admin_ID     string `param:"id_user" validate:"required"`
	Fullname     string `json:"fullname" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	Role         string `json:"role" `
	Phone        string `json:"phone" validate:"required"`
	Status       bool   `json:"status" `
	Verification bool   `json:"verification" `
}

type AdminDeleteRequest struct {
	Admin_ID string `param:"id_user" validate:"required"`
}
