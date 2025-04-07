package http

import "github.com/ilhamdcp/friendly-potato/internal/service"

type Handler struct {
	userService    service.UserService
	contactService service.ContactService
}

func NewHandler(userService service.UserService, contactService service.ContactService) *Handler {
	return &Handler{
		userService:    userService,
		contactService: contactService,
	}
}
