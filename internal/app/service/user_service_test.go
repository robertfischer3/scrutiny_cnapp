package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/myproject/internal/app/service"
)

func TestUserService_GetUser(t *testing.T) {
    // Setup mock repositories, etc.
    
    // Create service with mocks
    userService := service.NewUserService(/* pass mocks */)
    
    // Run test cases
    t.Run("Should return user by ID", func(t *testing.T) {
        // Test implementation
        user, err := userService.GetUserByID(1)
        
        // Assertions
        assert.NoError(t, err)
        assert.Equal(t, 1, user.ID)
        assert.Equal(t, "John Doe", user.Name)
    })
}