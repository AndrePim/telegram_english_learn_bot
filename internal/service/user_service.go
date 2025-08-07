package service

import (
	"github.com/AndrePim/telegram_english_learn_bot/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser регистрирует или обновляет пользователя
func (s *UserService) RegisterUser(userID int64, username, firstName, lastName string) error {
	user := &repository.User{
		ID:        userID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		State:     "idle",
	}

	return s.userRepo.CreateOrUpdateUser(user)
}

// GetUser получает пользователя
func (s *UserService) GetUser(userID int64) (*repository.User, error) {
	return s.userRepo.GetUser(userID)
}

// UpdateUserState обновляет состояние пользователя
func (s *UserService) UpdateUserState(userID int64, state string) error {
	return s.userRepo.UpdateUserState(userID, state)
}
