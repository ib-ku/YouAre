package usecase

import (
	"log"
	"time"
	"user-service/internal/entity"
	"user-service/internal/repository"

	gmail "user-service/internal/email"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string) (*entity.User, error)
	Login(email, password string) (*entity.TokenPair, error)
}

type authUsecase struct {
	repo         repository.UserRepository
	jwtSecret    string
	accessExpire time.Duration
}

func NewAuthUsecase(repo repository.UserRepository, jwtSecret string, accessExpire time.Duration) *authUsecase {
	return &authUsecase{
		repo:         repo,
		jwtSecret:    jwtSecret,
		accessExpire: accessExpire,
	}
}

// realization
func (a *authUsecase) Register(email, password string) (*entity.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = a.repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	sender := gmail.NewEmailSender()
	err = sender.Send(email, "Welcome!", "<h1>Thanks for registration</h1>")
	if err != nil {
		log.Fatalf("failed to send welcome email: %v", err)
	}
	log.Println("user registrated")
	log.Println("email sent")

	return newUser, nil
}

func (a *authUsecase) Login(email, password string) (*entity.TokenPair, error) {
	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, entity.ErrInvalidCredentials
	}

	accessToken, _, err := a.generateJWT(*user)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken: accessToken,
	}, nil
}

// utils
func (a *authUsecase) generateJWT(user entity.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(a.accessExpire)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.Hex(),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
