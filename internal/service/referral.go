package service

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"referral-system-test/internal/model"
	"referral-system-test/internal/repository"
	"time"
)

const ReferralCodeLength = 8

type UserService interface {
	RegisterUser(email, password string, referrerID *uint) error
	Login(email, password string) (string, error)
	CreateReferralCode(userID uint, code string, expiryDate time.Time) error
	GetReferralCode(email string) (*model.ReferralCode, error)
	GetReferrals(referrerID uint) ([]model.User, error)
	GenerateReferralCode(userID uint, expireDate time.Time) string
	DeleteReferralCode(userID uint) error
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *userService) RegisterUser(email, password string, referrerID *uint) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Email:      email,
		Password:   string(hashedPassword),
		ReferrerID: referrerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.CreateUser(user)
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
	})
	return token.SignedString(s.jwtSecret)
}

func (s *userService) CreateReferralCode(userID uint, code string, expiryDate time.Time) error {
	referralCode := &model.ReferralCode{
		Code:       code,
		UserID:     userID,
		ExpiryDate: expiryDate,
		CreatedAt:  time.Now(),
	}
	return s.repo.CreateReferralCode(referralCode)
}

func (s *userService) GetReferralCode(email string) (*model.ReferralCode, error) {
	return s.repo.FindReferralCodeByEmail(email)
}

func (s *userService) GetReferrals(referrerID uint) ([]model.User, error) {
	return s.repo.GetReferrals(referrerID)
}

func (s *userService) GenerateReferralCode(userID uint, expireDate time.Time) string {
	code := generateRandomString(ReferralCodeLength)

	s.repo.SaveReferralCode(userID, code, expireDate)

	return code
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Обработка ошибки генерации случайных байтов
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}

func (s *userService) DeleteReferralCode(userID uint) error {
	return s.repo.DeleteReferralCode(userID)
}
