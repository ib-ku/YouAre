package repository

import "user-service/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetUserById(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Delete(id string) error
}
