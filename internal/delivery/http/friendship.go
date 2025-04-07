package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
)

func (h *Handler) AddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var friendship domain.Contact
	err := json.NewDecoder(r.Body).Decode(&friendship)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if friendship.Username == "" || friendship.ContactUsername == "" {
		h.ErrorJSON(w, "Username and contactUsername are required", http.StatusBadRequest)
		return
	}

	data, err := h.contactService.AddContact(r.Context(), &friendship)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to add friend: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, data)
}

func (h *Handler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var friendship domain.Contact
	err := json.NewDecoder(r.Body).Decode(&friendship)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if friendship.Username == "" || friendship.ContactUsername == "" {
		h.ErrorJSON(w, "Username and contactUsernname are required", http.StatusBadRequest)
		return
	}

	err = h.contactService.RemoveContact(r.Context(), &friendship)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, "Contact removed successfully")
}

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		h.ErrorJSON(w, "Username is required", http.StatusBadRequest)
		return
	}

	contacts, err := h.contactService.GetContacts(r.Context(), username)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to get friends: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, contacts)
}
