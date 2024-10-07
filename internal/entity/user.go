package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId             uuid.UUID `json:"id_user" gorm:"column:id_user"`
	Fullname           string    `json:"fullname"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	Phone              string    `json:"phone"`
	Role               string    `json:"role"`
	Status             bool      `json:"status"`
	ResetCode          string    `json:"reset_code"`
	ResetCodeExpiresAt time.Time `json:"reset_code_expires_at"`
	Auditable
	Verification     bool   `json:"verification"`
	VerificationCode string `json:"verification_code"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Role == "" {
		u.Role = "user"
	}
	if !u.Status {
		u.Status = true
	}
	return
}

func NewUser(fullname, email, password, phone, role string, status, verification bool) *User {
	return &User{
		UserId:       uuid.New(),
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Phone:        phone,
		Role:         role,
		Status:       status,
		Verification: verification,
		Auditable:    NewAuditable(),
	}
}

func UpdateUser(id_user uuid.UUID, fullname, email, password, phone, role string, status, verification bool) *User {
	return &User{
		UserId:       id_user,
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Phone:        phone,
		Role:         role,
		Status:       status,
		Verification: verification,
		Auditable:    UpdateAuditable(),
	}
}
