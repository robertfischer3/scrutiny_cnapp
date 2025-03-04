package service

import (
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// FindByID mocks the FindByID method of the UserRepository interface
func (m *MockUserRepository) FindByID(id int) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

// FindAll mocks the FindAll method of the UserRepository interface
func (m *MockUserRepository) FindAll() ([]User, error) {
	args := m.Called()
	return args.Get(0).([]User), args.Error(1)
}

// Create mocks the Create method of the UserRepository interface
func (m *MockUserRepository) Create(user User) (User, error) {
	args := m.Called(user)
	return args.Get(0).(User), args.Error(1)
}

// Update mocks the Update method of the UserRepository interface
func (m *MockUserRepository) Update(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Delete mocks the Delete method of the UserRepository interface
func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}