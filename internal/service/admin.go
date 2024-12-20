package service

import (
	"errors" // Import log package

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/email"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/encrypt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	LoginAdmin(email string, password string) (string, error)
	FindAllUser(page int) ([]entity.User, error)
	CreateAdmin(admin *entity.Admin) (*entity.Admin, error)
	UpdateAdmin(admin *entity.Admin) (*entity.Admin, error)
	DeleteAdmin(admin uuid.UUID) (bool, error)
	EmailExists(email string) bool
	CheckUserExists(id uuid.UUID) (bool, error)
}

type adminService struct {
	adminRepository repository.AdminRepository
	encryptTool     encrypt.EncryptTool
	emailSender     *email.EmailSender
}

func NewAdminService(adminRepository repository.AdminRepository,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender) *adminService {
	return &adminService{
		adminRepository: adminRepository,
		encryptTool:     encryptTool,
		emailSender:     emailSender,
	}
}

func (s *adminService) LoginAdmin(email string, password string) (string, error) {
	admin, err := s.adminRepository.FindAdminByEmail(email)
	if err != nil {
		return "", errors.New("wrong input email/password")
	}
	if admin.Role != "admin" {
		return "", errors.New("you dont have access")
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	// Dekripsi nomor telepon jika perlu
	admin.Phone, _ = s.encryptTool.Decrypt(admin.Phone)
	return "Login successful", nil
}

func (s *adminService) FindAllUser(page int) ([]entity.User, error) {
	admin, err := s.adminRepository.FindAllUser(page)
	if err != nil {
		return nil, err
	}

	formattedAdmin := make([]entity.User, 0)
	for _, v := range admin {
		v.Phone, _ = s.encryptTool.Decrypt(v.Phone)
		formattedAdmin = append(formattedAdmin, v)
	}

	return formattedAdmin, nil
}

func (s *adminService) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin.Password = string(hashedPassword)

	newAdmin, err := s.adminRepository.CreateAdmin(admin)
	if err != nil {
		return nil, err
	}
	emailAddr := newAdmin.Email
	err = s.emailSender.SendWelcomeEmail(emailAddr, newAdmin.Fullname, "")

	if err != nil {
		return nil, err
	}

	resetCode := generateResetCode()
	err = s.emailSender.SendVerificationEmail(newAdmin.Email, newAdmin.Fullname, resetCode)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *adminService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.adminRepository.CheckUserExists(id)
}

func (s *adminService) UpdateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.Password = string(hashedPassword)
	}
	if admin.Phone != "" {
		admin.Phone, _ = s.encryptTool.Encrypt(admin.Phone)
	}

	return s.adminRepository.UpdateAdmin(admin)
}

func (s *adminService) DeleteAdmin(id_user uuid.UUID) (bool, error) {
	user, err := s.adminRepository.FindAdminByID(id_user)
	if err != nil {
		return false, err
	}

	return s.adminRepository.DeleteAdmin(user)
}

func (s *adminService) EmailExists(email string) bool {
	_, err := s.adminRepository.FindAdminByEmail(email)
	if err != nil {
		return false
	}
	return true
}
