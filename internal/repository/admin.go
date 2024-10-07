package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (*entity.Admin, error)
	FindAdminByID(id_user uuid.UUID) (*entity.Admin, error)
	FindByRole(role string, users *[]entity.User) error
	FindAllUser(page int) ([]entity.User, error)
	CreateAdmin(admin *entity.Admin) (*entity.Admin, error)
	UpdateAdmin(admin *entity.Admin) (*entity.Admin, error)
	DeleteAdmin(admin *entity.Admin) (bool, error)
	SaveVerifCode(userID uuid.UUID, resetCode string) error
	CheckUserExists(id uuid.UUID) (bool, error)
}

type adminRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewAdminRepository(db *gorm.DB, cacheable cache.Cacheable) *adminRepository {
	return &adminRepository{db: db, cacheable: cacheable}
}

func (r *adminRepository) FindAdminByID(id_user uuid.UUID) (*entity.Admin, error) {
	admin := new(entity.Admin)
	if err := r.db.Where("id_user = ?", id_user).Take(admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}

func (r *adminRepository) FindAdminByEmail(email string) (*entity.Admin, error) {
	admin := new(entity.Admin)
	if err := r.db.Where("email = ?", email).Take(admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}

func (r *adminRepository) FindByRole(role string, users *[]entity.User) error {
	return r.db.Where("role = ?", role).Find(users).Error
}

func (r *adminRepository) FindAllUser(page int) ([]entity.User, error) {
	var users []entity.User
	key := fmt.Sprintf("FindAllUsers_page_%d", page)
	const pageSize = 10

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
			return users, err
		}
		marshalledusers, _ := json.Marshal(users)
		err := r.cacheable.Set(key, marshalledusers, 5*time.Minute)
		if err != nil {
			return users, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			return users, err
		}
	}
	return users, nil
}

func (r *adminRepository) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if err := r.db.Create(&admin).Error; err != nil {
		return admin, err
	}
	r.cacheable.Delete("FindAllUsers_page_1")
	return admin, nil

}

func (r *adminRepository) CheckUserExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Admin{}).Where("id_user = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *adminRepository) UpdateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if admin.Email != "" {
		fields["email"] = admin.Email
	}
	if admin.Password != "" {
		fields["password"] = admin.Password
	}
	if admin.Role != "" {
		fields["role"] = admin.Role
	}

	// Update the database in one query.
	if err := r.db.Model(admin).Where("id_user = ?", admin.UserId).Updates(fields).Error; err != nil {
		return admin, err
	}
	r.cacheable.Delete("FindAllUsers_page_1")
	return admin, nil
}

func (r *adminRepository) DeleteAdmin(admin *entity.Admin) (bool, error) {
	if err := r.db.Delete(&entity.Admin{}, admin.UserId).Error; err != nil {
		return false, err
	}
	r.cacheable.Delete("FindAllUsers_page_1")
	return true, nil
}

func (r *adminRepository) SaveVerifCode(id_user uuid.UUID, resetCode string) error {
	return r.db.Model(&entity.User{}).Where("id_user = ?", id_user).Updates(map[string]interface{}{
		"verification_code": resetCode,
	}).Error
}
