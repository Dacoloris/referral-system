package repository

import (
	"referral-system-test/internal/model"
	"referral-system-test/pkg/db"
	"time"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	CreateReferralCode(code *model.ReferralCode) error
	FindReferralCodeByEmail(email string) (*model.ReferralCode, error)
	GetReferrals(referrerID uint) ([]model.User, error)
	SaveReferralCode(userID uint, code string, expireDate time.Time) error
	DeleteReferralCode(userID uint) error
}

type userRepository struct {
	db *db.Db
}

func NewUserRepository(db *db.Db) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) CreateReferralCode(code *model.ReferralCode) error {
	return r.db.Create(code).Error
}

func (r *userRepository) FindReferralCodeByEmail(email string) (*model.ReferralCode, error) {
	var code model.ReferralCode
	err := r.db.Where("user_id = (SELECT id FROM users WHERE email = ?)", email).First(&code).Error
	return &code, err
}

func (r *userRepository) GetReferrals(referrerID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("referrer_id = ?", referrerID).Find(&users).Error
	return users, err
}

// DeleteReferralCode - метод для удаления реферального кода пользователя
func (r *userRepository) DeleteReferralCode(userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&model.ReferralCode{}).Error; err != nil {
		return err
	}
	return nil
}

// SaveReferralCode - метод для сохранения реферального кода в базе данных
func (r *userRepository) SaveReferralCode(userID uint, code string, expireDate time.Time) error {
	referral := model.ReferralCode{
		UserID:     userID,
		Code:       code,
		ExpiryDate: expireDate,
	}

	// Проверяем, существует ли уже реферальный код для данного пользователя
	var existingReferral model.ReferralCode
	if err := r.db.Where("user_id = ?", userID).First(&existingReferral).Error; err == nil {
		// Если код уже существует, обновляем его
		existingReferral.Code = code
		return r.db.Save(&existingReferral).Error
	}

	// Если кода нет, создаем новый
	return r.db.Create(&referral).Error
}
