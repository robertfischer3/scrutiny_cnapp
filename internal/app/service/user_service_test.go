package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockUserRepository)
	
	// Create test user
	testUser := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Role:     "user",
		Active:   true,
		CreatedAt: "2025-03-01T10:00:00Z",
		UpdatedAt: "2025-03-01T10:00:00Z",
	}
	
	// Setup expectations
	mockRepo.On("FindByID", 1).Return(testUser, nil)
	
	// Create service with mock
	userService := NewUserService(mockRepo)
	
	// Run test cases
	t.Run("Should return user by ID", func(t *testing.T) {
		// Test implementation
		user, err := userService.GetUserByID(1)
		
		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		
		// Verify that our expectations were met
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Should handle invalid user ID", func(t *testing.T) {
		// Test with invalid ID
		user, err := userService.GetUserByID(-1)
		
		// Assertions
		assert.Error(t, err)
		assert.Equal(t, User{}, user)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockUserRepository)
	
	// Create test users
	testUsers := []User{
		{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			Role:     "user",
			Active:   true,
		},
		{
			ID:       2,
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Role:     "admin",
			Active:   true,
		},
	}
	
	// Setup expectations
	mockRepo.On("FindAll").Return(testUsers, nil)
	
	// Create service with mock
	userService := NewUserService(mockRepo)
	
	// Run test
	t.Run("Should return all users", func(t *testing.T) {
		// Test implementation
		users, err := userService.GetAllUsers()
		
		// Assertions
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, "Jane Smith", users[1].Name)
		
		// Verify that our expectations were met
		mockRepo.AssertExpectations(t)
	})
}