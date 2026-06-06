package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Flow_gym_go_project/repository"
	"Flow_gym_go_project/services"
)

// RegisterRequest represents the JSON payload required to create a new user.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration requests.
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only POST requests are allowed for registration.
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields.
		if req.Name == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "name, email and password are required", http.StatusBadRequest)
			return
		}

		// Hash the password before storing it in the database.
		passwordHash, err := services.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "could not hash password", http.StatusInternalServerError)
			return
		}

		// Create the user repository.
		userRepo := repository.NewUserRepository(db)

		// Create a new user with the default role "user".
		user, err := userRepo.Create(req.Name, req.Email, passwordHash, "user")
		if err != nil {
			http.Error(w, "could not create user", http.StatusInternalServerError)
			return
		}

		// Return the created user as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// LoginRequest represents the JSON payload required to log in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles user login requests.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only POST requests are allowed for login.
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields.
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		// Create the user repository.
		userRepo := repository.NewUserRepository(db)

		// Look up the user by email.
		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {

			// Avoid revealing whether the email exists.
			if err == sql.ErrNoRows {
				http.Error(w, "invalid email or password", http.StatusUnauthorized)
				return
			}

			http.Error(w, "could not login user", http.StatusInternalServerError)
			return
		}

		// Compare the provided password with the stored hash.
		if !services.CheckPassword(req.Password, user.PasswordHash) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		// Store the authenticated user ID in an HTTP-only cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    strconv.Itoa(user.ID),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		// Return the authenticated user as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// LogoutHandler clears the authentication cookie.
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only POST requests are allowed for logout.
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Expire the user_id cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		// Return a logout confirmation.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "logged out successfully",
		})
	}
}

// MeHandler returns the currently authenticated user.
func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the authentication cookie.
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Convert the cookie value to a user ID.
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}

		// Create the user repository.
		userRepo := repository.NewUserRepository(db)

		// Retrieve the current user from the database.
		user, err := userRepo.GetByID(userID)
		if err != nil {

			// If the user no longer exists, the session is invalid.
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			http.Error(w, "could not fetch current user", http.StatusInternalServerError)
			return
		}

		// Return the current user as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}