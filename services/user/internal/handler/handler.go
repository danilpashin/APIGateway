package handler

import (
	"apigateway/services/user/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (u *UserHandler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Basic handler check"))
}
