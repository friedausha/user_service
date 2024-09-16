package handler

import (
	"context"
	"encoding/json"
	_ "errors"
	"fmt"
	"git.garena.com/frieda.hasanah/user_service/internal/dto"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"strings"
)

type Handler struct {
	authService model.IAuthService
	userService model.IUserService
}

// NewHandler returns a new handler instance
func NewHandler(authService model.IAuthService, userService model.IUserService) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
	}
}

// HandleRequest handles the raw TCP requests
func (h *Handler) HandleRequest(request string) string {
	requestParts := strings.SplitN(request, " ", 2)
	if len(requestParts) < 2 {
		return h.sendError("Invalid request format")
	}

	switch requestParts[0] {
	case "REGISTER":
		return h.Register(requestParts[1])
	case "LOGIN":
		return h.Login(requestParts[1])
	case "VERIFY_TOKEN":
		return h.VerifyToken(requestParts[1])
	default:
		return h.sendError("Unknown command")
	}
}

func (h *Handler) sendError(message string) string {
	return fmt.Sprintf(`{"error": "%s"}`, message)
}

func (h *Handler) Register(payload string) string {
	var u model.User
	if err := json.Unmarshal([]byte(payload), &u); err != nil {
		return h.sendError("Invalid request payload")
	}

	u, err := h.userService.Register(context.Background(), u)
	if err != nil {
		return h.sendError(err.Error())
	}

	response, _ := json.Marshal(dto.TransformRegisterResponse(u))
	return string(response)
}

func (h *Handler) Login(payload string) string {
	var u model.User
	if err := json.Unmarshal([]byte(payload), &u); err != nil {
		return h.sendError("Invalid request payload")
	}

	u, expIn, err := h.authService.Login(context.Background(), u.Username, u.Password)
	if err != nil {
		return h.sendError(err.Error())
	}

	response, _ := json.Marshal(dto.TransformLoginResponse(u, expIn))
	return string(response)
}

func (h *Handler) VerifyToken(token string) string {
	token = strings.TrimPrefix(token, "Bearer ")
	if _, err := h.authService.VerifyToken(token); err != nil {
		return h.sendError(err.Error())
	}

	response, _ := json.Marshal(dto.TransformVerifyResponse())
	return string(response)
}
