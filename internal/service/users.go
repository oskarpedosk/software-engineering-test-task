package service

import (
	"cruder/internal/model"
	"cruder/internal/repository"
	"fmt"
	"net/mail"
	"strings"
)

type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	GetUserByUUID(uuid string) (*model.User, error)
	CreateUser(request *model.CreateUserRequest) (*model.User, error)
	UpdateUser(uuid string, request *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(uuid string) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (userService *userService) GetAllUsers() ([]model.User, error) {
	return userService.userRepository.GetAll()
}

func (userService *userService) validateUserExists(user *model.User, err error) (*model.User, error) {
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}

func (userService *userService) updateFieldIfProvided(target *string, value string) {
	if value != "" {
		*target = value
	}
}

func (userService *userService) GetUserByUsername(username string) (*model.User, error) {
	user, err := userService.userRepository.GetByUsername(username)
	return userService.validateUserExists(user, err)
}

func (userService *userService) GetUserByID(id int64) (*model.User, error) {
	user, err := userService.userRepository.GetByID(id)
	return userService.validateUserExists(user, err)
}

func (userService *userService) GetUserByUUID(uuid string) (*model.User, error) {
	if err := userService.validateNonEmpty(uuid); err != nil {
		return nil, err
	}

	user, err := userService.userRepository.GetByUUID(uuid)
	return userService.validateUserExists(user, err)
}

func (userService *userService) CreateUser(request *model.CreateUserRequest) (*model.User, error) {
	if err := userService.validateCreateRequest(request); err != nil {
		return nil, err
	}

	user := &model.User{
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
	}

	if err := userService.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *userService) UpdateUser(uuid string, request *model.UpdateUserRequest) (*model.User, error) {
	if err := userService.validateNonEmpty(uuid); err != nil {
		return nil, err
	}

	existing, err := userService.validateUserExists(userService.userRepository.GetByUUID(uuid))
	if err != nil {
		return nil, err
	}

	if err := userService.validateUpdateRequest(request); err != nil {
		return nil, err
	}

	userService.updateFieldIfProvided(&existing.Username, request.Username)
	userService.updateFieldIfProvided(&existing.Email, request.Email)
	userService.updateFieldIfProvided(&existing.FullName, request.FullName)

	if err := userService.userRepository.Update(uuid, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (userService *userService) DeleteUser(uuid string) error {
	if err := userService.validateNonEmpty(uuid); err != nil {
		return err
	}

	if _, err := userService.validateUserExists(userService.userRepository.GetByUUID(uuid)); err != nil {
		return err
	}

	return userService.userRepository.Delete(uuid)
}

func (userService *userService) validateRequiredField(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s: %w", fieldName, model.ErrEmptyField)
	}
	return nil
}

func (userService *userService) validateCreateRequest(request *model.CreateUserRequest) error {
	if err := userService.validateRequiredField("username", request.Username); err != nil {
		return err
	}

	if err := userService.validateRequiredField("email", request.Email); err != nil {
		return err
	}

	if !isValidEmail(request.Email) {
		return model.ErrInvalidEmail
	}

	return nil
}

func (userService *userService) validateUpdateRequest(request *model.UpdateUserRequest) error {
	if strings.TrimSpace(request.Email) != "" && !isValidEmail(request.Email) {
		return model.ErrInvalidEmail
	}
	return nil
}

func (userService *userService) validateNonEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return model.ErrEmptyField
	}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
