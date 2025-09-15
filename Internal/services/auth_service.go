package services

import (
	"time"

	"github.com/Archnick/go-ecommerce/Internal/models"
	"github.com/Archnick/go-ecommerce/Internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) UpdateUser(user *models.User, payload models.UpdateUserPayload) error {
	return s.userRepo.UpdateWithPayload(user, payload)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) RegisterUser(payload models.UserPayload) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
		Role:     string(models.CustomerRole)}

	err = s.userRepo.Create(user)
	return user, err
}

func (s *UserService) AuthorizeUser(payload models.UserPayload) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(payload.Email)
	if nil != err {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) SetRefreshToken(user *models.User, refreshToken string) error {
	expiry := time.Now().Add(7 * 24 * time.Hour)
	if refreshToken == "" {
		expiry = time.Time{}
	}

	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = expiry
	return s.userRepo.Update(user)
}
