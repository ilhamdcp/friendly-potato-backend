package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
	"github.com/ilhamdcp/friendly-potato/internal/dto"
)

func (h *Handler) ReturnJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&dto.Response{Status: http.StatusOK, Data: data})
}

func (h *Handler) ErrorJSON(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&dto.Response{Status: code, Data: err})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUser, err := h.userService.CreateUser(r.Context(), &user)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}
	h.ReturnJSON(w, newUser)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Path[len("/users/"):] // Extract user ID from path
	if userID == "" {
		h.ErrorJSON(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if valid := h.userService.AuthenticateUser(r.Context(), r.Header.Get("Authorization")); !valid {
		h.ErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
		return
	}
	if user == nil {
		h.ReturnJSON(w, errors.New("user not found"))
		return
	}

	h.ReturnJSON(w, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Path[len("/users/"):] // Extract user ID from path
	if userID == "" {
		h.ErrorJSON(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.ID != userID {
		h.ErrorJSON(w, "User ID in path does not match ID in body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUser(r.Context(), &user)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, true)
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (h *Handler) SignInUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userService.SignInUser(r.Context(), &user)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to sign in user: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, token)
}
func (h *Handler) SignOutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success, err := h.userService.SignOutUser(r.Context(), user.Username)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to sign out user: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, success)
}

func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		h.ErrorJSON(w, "Token is required", http.StatusBadRequest)
		return
	}

	valid := h.userService.AuthenticateUser(r.Context(), token)
	h.ReturnJSON(w, valid)
}
