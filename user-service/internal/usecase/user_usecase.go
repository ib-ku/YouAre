package usecase

import (
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/cache"
	"user-service/internal/entity"
	"user-service/internal/repository"
)

// user usecase
type UserUsecase interface {
	GetUserById(ID string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	DeleteUser(ID string) error
}

type userUsecase struct {
	repo  repository.UserRepository
	cache *cache.RedisCache
}

func NewUserUsecase(repo repository.UserRepository, c *cache.RedisCache) *userUsecase {
	return &userUsecase{
		repo:  repo,
		cache: c,
	}
}

// realization
func (u *userUsecase) GetUserById(ID string) (*entity.User, error) {
	if ID == "" {
		return nil, entity.ErrInvalidInput
	}

	cacheKey := "user:" + ID

	cached, err := u.cache.Get(cacheKey)
	if err == nil {
		fmt.Println("CACHE HIT for", cacheKey)
		var user entity.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	fmt.Println("CACHE MISS for", cacheKey)
	user, err := u.repo.GetUserById(ID)
	if err != nil {
		return nil, err
	}

	userBytes, err := json.Marshal(user)
	if err == nil {
		u.cache.Set(cacheKey, string(userBytes), 10*time.Minute)
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

func (u *userUsecase) GetAll() ([]*entity.User, error) {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecase) DeleteUser(ID string) error {
	if ID == "" {
		return entity.ErrInvalidInput
	}

	_, err := u.repo.GetUserById(ID)
	if err != nil {
		return err
	}

	u.cache.Delete("user:" + ID)
	return u.repo.Delete(ID)
}
