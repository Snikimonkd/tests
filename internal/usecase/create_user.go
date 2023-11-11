package usecase

import (
	"lab1/internal/model"
)

//go:generate mockgen -destination=./mocks/mock_CreateUserRepository.go -package=mocks lab1/internal/usecase CreateUserRepository
type CreateUserRepository interface {
	CheckUserExist(email string) (bool, error)
	CreateUser(user model.User) error
}

func (u usecase) CreateUser(user model.User) error {
	exist, err := u.createUserRepository.CheckUserExist(user.Email)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	err = u.createUserRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
