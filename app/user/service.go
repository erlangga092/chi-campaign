package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	IsEmailAvailable(input CheckEmailAvailableInput) (bool, error)
	LoginUser(input LoginUserInput) (User, error)
	GetUserByID(userID int) (User, error)
	UploadAvatar(userID int, fileLocation string) (User, error)
}

type userService struct {
	userRepository Repository
}

func NewUserService(userRepository Repository) Service {
	return &userService{userRepository}
}

func (s *userService) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.userRepository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) IsEmailAvailable(input CheckEmailAvailableInput) (bool, error) {
	user, err := s.userRepository.FindByEmail(input.Email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) LoginUser(input LoginUserInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("email or password not match")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, errors.New("email or password not match")
	}

	return user, nil
}

func (s *userService) GetUserByID(userID int) (User, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) UploadAvatar(userID int, fileLocation string) (User, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.userRepository.Update(user.ID, user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}
