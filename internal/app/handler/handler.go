package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/robertfischer3/scrutiny_cnapp/internal/app/service"
	"github.com/robertfischer3/scrutiny_cnapp/internal/pkg/logger"
)

// UserHandler handles HTTP requests for user resources
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles GET requests for a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers handles GET requests for all users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// LoggingMiddleware logs information about each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log the request
		log := logger.GetLogger()
		log.WithFields(map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"duration":   time.Since(start),
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("Request handled")
	})
}

// RegisterHandlers registers all HTTP handlers to the router
func RegisterHandlers(r *mux.Router) {
	// For now, we'll create a simple health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Add API version prefix
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	
	// Health routes
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	}).Methods("GET")
	
	// TODO: Add actual user handlers when repositories are implemented
	// Example of how you would set up user routes:
	/*
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)
	
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	*/
}