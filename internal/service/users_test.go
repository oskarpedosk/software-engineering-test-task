package service

import (
	"cruder/internal/model"
	"cruder/internal/repository"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users      map[string]*model.User
	nextID     int
	shouldFail bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*model.User),
		nextID: 1,
	}
}

func (userRepository *mockUserRepository) GetAll() ([]model.User, error) {
	if userRepository.shouldFail {
		return nil, assert.AnError
	}

	users := make([]model.User, 0, len(userRepository.users))
	for _, user := range userRepository.users {
		users = append(users, *user)
	}
	return users, nil
}

func (userRepository *mockUserRepository) GetByUsername(username string) (*model.User, error) {
	if userRepository.shouldFail {
		return nil, assert.AnError
	}

	for _, user := range userRepository.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (userRepository *mockUserRepository) GetByID(id int64) (*model.User, error) {
	if userRepository.shouldFail {
		return nil, assert.AnError
	}

	for _, user := range userRepository.users {
		if int64(user.ID) == id {
			return user, nil
		}
	}
	return nil, nil
}

func (userRepository *mockUserRepository) GetByUUID(uuid string) (*model.User, error) {
	if userRepository.shouldFail {
		return nil, assert.AnError
	}

	user, exists := userRepository.users[uuid]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (userRepository *mockUserRepository) Create(user *model.User) error {
	if userRepository.shouldFail {
		return assert.AnError
	}

	user.ID = userRepository.nextID
	user.UUID = generateMockUUID(userRepository.nextID)
	userRepository.nextID++
	userRepository.users[user.UUID] = user
	return nil
}

func (userRepository *mockUserRepository) Update(uuid string, user *model.User) error {
	if userRepository.shouldFail {
		return assert.AnError
	}

	existing, exists := userRepository.users[uuid]
	if !exists {
		return nil
	}

	existing.Username = user.Username
	existing.Email = user.Email
	existing.FullName = user.FullName
	return nil
}

func (userRepository *mockUserRepository) Delete(uuid string) error {
	if userRepository.shouldFail {
		return assert.AnError
	}

	delete(userRepository.users, uuid)
	return nil
}

func generateMockUUID(id int) string {
	return fmt.Sprintf("123e4567-e89b-12d3-a456-42661417%04d", id)
}

func setupTest() (*mockUserRepository, UserService) {
	mockRepo := newMockUserRepository()
	userService := NewUserService(mockRepo)
	return mockRepo, userService
}

func TestShouldGetAllUsers(t *testing.T) {
	// given
	repo, svc := setupTest()
	user1 := &model.User{ID: 1, UUID: generateMockUUID(1), Username: "user1", Email: "user1@test.com"}
	user2 := &model.User{ID: 2, UUID: generateMockUUID(2), Username: "user2", Email: "user2@test.com"}
	repo.users[user1.UUID] = user1
	repo.users[user2.UUID] = user2

	// when
	result, err := svc.GetAllUsers()

	// then
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestShouldGetUserByUsername(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	// when
	result, err := userService.GetUserByUsername("test_user")

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, "test_user", result.Username)
}

func TestShouldReturnErrorWhenUserNotFoundByUsername(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	result, err := userService.GetUserByUsername("nonexistent")

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrUserNotFound)
	assert.Nil(test, result)
}

func TestShouldGetUserByID(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	// when
	result, err := userService.GetUserByID(1)

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, 1, result.ID)
}

func TestShouldReturnErrorWhenUserNotFoundByID(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	result, err := userService.GetUserByID(999)

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrUserNotFound)
	assert.Nil(test, result)
}

func TestShouldGetUserByUUID(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	// when
	result, err := userService.GetUserByUUID("123e4567-e89b-12d3-a456-426614174000")

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, "123e4567-e89b-12d3-a456-426614174000", result.UUID)
}

func TestShouldReturnErrorWhenUserNotFoundByUUID(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	result, err := userService.GetUserByUUID("nonexistent-uuid")

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrUserNotFound)
	assert.Nil(test, result)
}

func TestShouldReturnErrorWhenUUIDIsEmpty(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	result, err := userService.GetUserByUUID("")

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrEmptyField)
	assert.Nil(test, result)
}

func TestShouldCreateUser(test *testing.T) {
	// given
	_, userService := setupTest()
	request := &model.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		FullName: "New User",
	}

	// when
	result, err := userService.CreateUser(request)

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, "newuser", result.Username)
	assert.NotEmpty(test, result.UUID)
}

func TestShouldReturnErrorWhenCreateUserWithInvalidData(test *testing.T) {
	tests := []struct {
		name    string
		request *model.CreateUserRequest
		wantErr error
	}{
		{
			name: "Empty username",
			request: &model.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				FullName: "Test",
			},
			wantErr: model.ErrEmptyField,
		},
		{
			name: "Empty email",
			request: &model.CreateUserRequest{
				Username: "test_user",
				Email:    "",
				FullName: "Test",
			},
			wantErr: model.ErrEmptyField,
		},
		{
			name: "Invalid email format",
			request: &model.CreateUserRequest{
				Username: "test_user",
				Email:    "invalid-email",
				FullName: "Test",
			},
			wantErr: model.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(subTest *testing.T) {
			// given
			_, userService := setupTest()

			// when
			result, err := userService.CreateUser(tt.request)

			// then
			assert.Error(subTest, err)
			assert.ErrorIs(subTest, err, tt.wantErr)
			assert.Nil(subTest, result)
		})
	}
}

func TestShouldUpdateUser(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "olduser",
		Email:    "old@example.com",
		FullName: "Old User",
	}
	mockRepo.users[user.UUID] = user

	request := &model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
	}

	// when
	result, err := userService.UpdateUser(user.UUID, request)

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, "updateduser", result.Username)
	assert.Equal(test, "updated@example.com", result.Email)
}

func TestShouldReturnErrorWhenUpdateNonExistentUser(test *testing.T) {
	// given
	_, userService := setupTest()

	request := &model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
	}

	// when
	result, err := userService.UpdateUser("nonexistent-uuid", request)

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrUserNotFound)
	assert.Nil(test, result)
}

func TestShouldReturnErrorWhenUpdateUserWithInvalidEmail(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	request := &model.UpdateUserRequest{
		Email: "invalid-email",
	}

	// when
	result, err := userService.UpdateUser(user.UUID, request)

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrInvalidEmail)
	assert.Nil(test, result)
}

func TestShouldPartiallyUpdateUser(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	request := &model.UpdateUserRequest{
		Email: "newemail@example.com",
	}

	// when
	result, err := userService.UpdateUser(user.UUID, request)

	// then
	assert.NoError(test, err)
	assert.NotNil(test, result)
	assert.Equal(test, "newemail@example.com", result.Email)
	assert.Equal(test, "test_user", result.Username)
}

func TestShouldDeleteUser(test *testing.T) {
	// given
	mockRepo, userService := setupTest()
	user := &model.User{
		ID:       1,
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Username: "test_user",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockRepo.users[user.UUID] = user

	// when
	err := userService.DeleteUser(user.UUID)

	// then
	assert.NoError(test, err)
	_, exists := mockRepo.users[user.UUID]
	assert.False(test, exists)
}

func TestShouldReturnErrorWhenDeleteNonExistentUser(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	err := userService.DeleteUser("nonexistent-uuid")

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrUserNotFound)
}

func TestShouldReturnErrorWhenDeleteUserWithEmptyUUID(test *testing.T) {
	// given
	_, userService := setupTest()

	// when
	err := userService.DeleteUser("")

	// then
	assert.Error(test, err)
	assert.ErrorIs(test, err, model.ErrEmptyField)
}

var _ repository.UserRepository = (*mockUserRepository)(nil)
