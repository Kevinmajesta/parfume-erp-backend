package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/internal/repository"
	"github.com/Kevinmajesta/webPemancingan/pkg/email"
	"github.com/Kevinmajesta/webPemancingan/pkg/encrypt"
	"github.com/Kevinmajesta/webPemancingan/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginUser(email string, password string) (string, error)
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id_user uuid.UUID) (bool, error)
	ResetPassword(resetCode string, newPassword string) error
	EmailExists(email string) bool
	GetUserProfileByID(userID string) (*entity.User, error)
	VerifUser(resetCode string) error
	CheckUserExists(id uuid.UUID) (bool, error)
	// GetCart(UserId uuid.UUID) (binder.GetCartResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	tokenUseCase   token.TokenUseCase
	encryptTool    encrypt.EncryptTool
	emailSender    *email.EmailSender
}

var InternalError = "internal server error"

func NewUserService(userRepository repository.UserRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender) *userService {
	return &userService{
		userRepository: userRepository,
		tokenUseCase:   tokenUseCase,
		encryptTool:    encryptTool,
		emailSender:    emailSender,
	}
}

func (s *userService) LoginUser(email string, password string) (string, error) {
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	// Lanjutkan dengan pembuatan token dan logika lainnya
	expiredTime := time.Now().Local().Add(24 * time.Hour)

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	// Dekripsi nomor telepon jika perlu
	user.Phone, _ = s.encryptTool.Decrypt(user.Phone)
	expiredTimeInJakarta := expiredTime.In(location)
	// Buat claims JWT
	claims := token.JwtCustomClaims{
		ID:    user.UserId.String(),
		Email: user.Email,
		Role:  "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Depublic",
			ExpiresAt: jwt.NewNumericDate(expiredTimeInJakarta),
		},
	}

	// Generate JWT token
	jwtToken, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("there is an error in the system")
	}

	return jwtToken, nil
}

func (s *userService) CreateUser(user *entity.User) (*entity.User, error) {
	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if user.Fullname == "" {
		return nil, errors.New("fullname cannot be empty")
	}
	if user.Phone == "" {
		return nil, errors.New("phone cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	emailAddr := newUser.Email
	err = s.emailSender.SendWelcomeEmail(emailAddr, newUser.Fullname, "")

	if err != nil {
		return nil, err
	}

	// resetCode := generateResetCode()
	// err = s.emailSender.SendVerificationEmail(newUser.Email, newUser.Fullname, resetCode)
	// if err != nil {
	// 	return nil, err
	// }
	// err = s.userRepository.SaveVerifCode(user.UserId, resetCode)
	// if err != nil {
	// 	return nil, err
	// }

	// notification := &entity.Notification{
	// 	UserID:  newUser.UserId,
	// 	Type:    "Registration",
	// 	Message: "User registration successful",
	// 	IsRead:  false,
	// }
	// err = s.notificationService.CreateNotification(notification)
	// if err != nil {
	// 	return nil, err
	// }

	return newUser, nil
}

func (s *userService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.userRepository.CheckUserExists(id)
}

func (s *userService) UpdateUser(user *entity.User) (*entity.User, error) {
	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if user.Fullname == "" {
		return nil, errors.New("fullname cannot be empty")
	}
	if user.Phone == "" {
		return nil, errors.New("phone cannot be empty")
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	updatedUser, err := s.userRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// notification := &entity.Notification{
	// 	UserID:  updatedUser.UserId,
	// 	Type:    "Update Profile",
	// 	Message: "Update Profile successful",
	// 	IsRead:  false,
	// }
	// err = s.notificationService.CreateNotification(notification)
	// if err != nil {
	// 	return nil, err
	// }

	return updatedUser, nil
}

func (s *userService) DeleteUser(id_user uuid.UUID) (bool, error) {
	user, err := s.userRepository.FindUserByID(id_user)
	if err != nil {
		return false, err
	}

	return s.userRepository.DeleteUser(user)
}

// func (s *userService) RequestPasswordReset(email string) error {
// 	user, err := s.userRepository.FindUserByEmail(email)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}

// 	resetCode := generateResetCode()
// 	expiresAt := time.Now().Add(1 * time.Hour)

// 	// err = s.userRepository.SaveResetCode(user.UserId, resetCode, expiresAt)
// 	// if err != nil {
// 	// 	return errors.New("failed to save reset code")
// 	// }

// 	// return s.emailSender.SendResetPasswordEmail(user.Email, user.Fullname, resetCode)
// }

func (s *userService) ResetPassword(resetCode string, newPassword string) error {
	user, err := s.userRepository.FindUserByResetCode(resetCode)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid reset code")
	}
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) EmailExists(email string) bool {
	_, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return false
	}
	return true
}

func (s *userService) GetUserProfileByID(userID string) (*entity.User, error) {
	// Konversi userID menjadi uuid.UUID
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Panggil metode dari userRepository untuk mencari profil pengguna berdasarkan ID
	user, err := s.userRepository.FindUserByID(userIDUUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) VerifUser(verifCode string) error {
	user, err := s.userRepository.FindUserByVerifCode(verifCode)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid verification code")
	}

	// Update verification status
	user.Verification = true

	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func generateResetCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}
