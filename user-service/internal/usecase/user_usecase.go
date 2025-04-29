package usecase

import (
	"user-service/internal/entity"
	"user-service/internal/repository"
)

// user usecase
type UserUsecase interface {
	GetUserById(ID string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetAll() []*entity.User
	DeleteUser(ID string) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *userUsecase {
	return &userUsecase{
		repo: repo,
	}
}

// realization
func (u *userUsecase) GetUserById(ID string) (*entity.User, error) {
	if ID == "" {
		return nil, entity.ErrInvalidInput
	}
	user, err := u.repo.GetUserById(ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetUserByEmail(email string) (*entity.User, error) {
	if email == "" {
		return nil, entity.ErrInvalidInput
	}
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetAll() []*entity.User {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil
	}

	return users
}

func (u *userUsecase) DeleteUser(ID string) error {
	if ID == "" {
		return entity.ErrInvalidInput
	}
	_, err := u.repo.GetUserById(ID)
	if err != nil {
		return err
	}
	return u.repo.Delete(ID)
}
