package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	UserId   uuid.UUID `gorm:"column:id_user;primary_key"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Phone    string    `json:"phone"`
	Auditable
	Verification     bool   `json:"verification"`
	VerificationCode string `json:"verification_code"`
}

func (u *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Role == "" {
		u.Role = "admin"
	}
	if u.Fullname == "" {
		u.Fullname = "admin"
	}
	return
}

func NewAdmin(fullname, email, password, role, phone string, verification bool) *Admin {
	return &Admin{
		UserId:       uuid.New(),
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    NewAuditable(),
		Verification: false,
	}
}

func UpdateAdmin(id_user uuid.UUID, fullname, email, password, role, phone string, verification bool) *Admin {
	return &Admin{
		UserId:       id_user,
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    UpdateAuditable(),
		Verification: false,
	}
}

func (Admin) TableName() string {
	return "users"
}
