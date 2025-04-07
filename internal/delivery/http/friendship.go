package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
)

func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var friendship domain.Friendship
	err := json.NewDecoder(r.Body).Decode(&friendship)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if friendship.UserID == "" || friendship.FriendID == "" {
		h.ErrorJSON(w, "User ID and Friend ID are required", http.StatusBadRequest)
		return
	}

	data, err := h.friendshipService.AddFriend(r.Context(), &friendship)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to add friend: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, data)
}

func (h *Handler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var friendship domain.Friendship
	err := json.NewDecoder(r.Body).Decode(&friendship)
	if err != nil {
		h.ErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if friendship.UserID == "" || friendship.FriendID == "" {
		h.ErrorJSON(w, "User ID and Friend ID are required", http.StatusBadRequest)
		return
	}

	err = h.friendshipService.RemoveFriend(r.Context(), &friendship)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, "Friend removed successfully")
}

func (h *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		h.ErrorJSON(w, "User ID is required", http.StatusBadRequest)
		return
	}

	friends, err := h.friendshipService.GetFriends(r.Context(), userID)
	if err != nil {
		h.ErrorJSON(w, fmt.Sprintf("Failed to get friends: %v", err), http.StatusInternalServerError)
		return
	}

	h.ReturnJSON(w, friends)
}
