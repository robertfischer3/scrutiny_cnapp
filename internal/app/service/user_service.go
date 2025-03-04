package service

import (
	"errors"
)

// User represents a user entity in the system
type User struct {
	ID       int
	Name     string
	Email    string
	Role     string
	Active   bool
	CreatedAt string
	UpdatedAt string
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id int) (User, error)
	FindAll() ([]User, error)
	Create(user User) (User, error)
	Update(user User) error
	Delete(id int) error
}

// UserService provides user-related operations
type UserService struct {
	repository UserRepository
}

// NewUserService creates a new UserService with the given repository
func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id int) (User, error) {
	if id <= 0 {
		return User{}, errors.New("invalid user ID")
	}
	
	return s.repository.FindByID(id)
}

// GetAllUsers retrieves all active users
func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repository.FindAll()
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user User) (User, error) {
	// Validate user data
	if user.Name == "" {
		return User{}, errors.New("user name cannot be empty")
	}
	if user.Email == "" {
		return User{}, errors.New("user email cannot be empty")
	}
	
	// Create the user
	return s.repository.Create(user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user User) error {
	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}
	
	// Ensure user exists
	_, err := s.repository.FindByID(user.ID)
	if err != nil {
		return err
	}
	
	return s.repository.Update(user)
}

// DeactivateUser deactivates a user
func (s *UserService) DeactivateUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}
	
	// Get the current user
	user, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	
	// Deactivate the user
	user.Active = false
	
	return s.repository.Update(user)
}