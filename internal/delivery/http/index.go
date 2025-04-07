package http

import "github.com/ilhamdcp/friendly-potato/internal/service"

type Handler struct {
	userService       service.UserService
	friendshipService service.FriendshipService
}

func NewHandler(userService service.UserService, friendshipService service.FriendshipService) *Handler {
	return &Handler{
		userService:       userService,
		friendshipService: friendshipService,
	}
}
