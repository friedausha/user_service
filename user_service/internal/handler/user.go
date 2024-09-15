package handler

import (
	"encoding/json"
	"git.garena.com/frieda.hasanah/user_service/internal/dto"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
)

type handler struct {
	authService model.IAuthService
	userService model.IUserService
}

// Init will initialize the REST handler for user service
func Init(e *echo.Group, authService model.IAuthService, userService model.IUserService) {
	h := handler{
		authService: authService,
		userService: userService,
	}

	e.POST("/api/register", h.Register)
	e.POST("/api/token", h.Login)
	e.POST("/api/token/verify", h.VerifyToken)
}

func (h *handler) sendError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

func (h handler) Register(c echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var u model.User
	if err = json.Unmarshal(b, &u); err != nil {
		return h.sendError(c, http.StatusBadRequest, "Invalid request payload")
	}

	if u, err = h.userService.Register(c.Request().Context(), u); err != nil {
		return h.sendError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dto.TransformRegisterResponse(u))
}

func (h handler) Login(c echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.sendError(c, http.StatusBadRequest, " Failed to read request body")
	}

	var u model.User
	if err = json.Unmarshal(b, &u); err != nil {
		return h.sendError(c, http.StatusBadRequest, "Invalid request payload")
	}

	var expIn int64
	if u, expIn, err = h.authService.Login(c.Request().Context(), u.Username, u.Password); err != nil {
		return h.sendError(c, http.StatusUnauthorized, err.Error())
	}

	response := dto.TransformLoginResponse(u, expIn)
	return c.JSON(http.StatusOK, response)
}

func (h handler) VerifyToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return h.sendError(c, http.StatusUnauthorized, "Missing Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return h.sendError(c, http.StatusUnauthorized, "Invalid token format")
	}

	if _, err := h.authService.VerifyToken(tokenStr); err != nil {
		return h.sendError(c, http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, dto.TransformVerifyResponse())
}
