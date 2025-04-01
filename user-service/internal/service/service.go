package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	customError "social-network/user-service/internal/errors"
	"social-network/user-service/internal/logger"
	"social-network/user-service/internal/repository"
)

type UserServiceInterface interface {
	Register(user *repository.User) (string, error)
	Login(user *repository.User) (string, error)
	GetUserProfile(login string) (*repository.User, error)
	UpdateUserProfile(login string, user *repository.User) error
}

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Register(user *repository.User) (string, error) {
	err := us.userRepository.RegisterUser(user)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	token, err := createJWTToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT token: %w", err)
	}
	return token, nil
}

func (us *UserService) Login(user *repository.User) (string, error) {
	err := us.userRepository.LoginUser(user)
	if err != nil {
		return "", fmt.Errorf("failed to login user: %w", err)
	}

	token, err := createJWTToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT token: %w", err)
	}
	return token, nil
}

func (us *UserService) GetUserProfile(login string) (*repository.User, error) {
	return us.userRepository.GetUserByLogin(login)
}

func (us *UserService) UpdateUserProfile(login string, user *repository.User) error {
	if user.Password != "" || user.Login != "" {
		return &customError.UpdateCredentialsError{}
	}
	return us.userRepository.UpdateUserProfile(login, user)
}

func createJWTToken(user *repository.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":    user.Name,
		"login":   user.Login,
		"user-id": user.Id,
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logger.Error(fmt.Sprintf("error creating jtw token: %v", err))
		return "", err
	}
	return token, nil
}
