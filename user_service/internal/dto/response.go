package dto

import "git.garena.com/frieda.hasanah/user_service/internal/model"

type RegisterResponseBody struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type LoginResponseBody struct {
	Message   string `json:"message"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

type VerifyResponseBody struct {
	Message string `json:"message"`
	Valid   bool   `json:"valid"`
}

func TransformLoginResponse(u model.User, expIn int64) LoginResponseBody {
	return LoginResponseBody{
		Message:   "User has logged in successfully",
		Token:     *u.Token,
		ExpiresIn: expIn,
	}
}

func TransformVerifyResponse() VerifyResponseBody {
	return VerifyResponseBody{
		Message: "User has been verified successfully",
		Valid:   true,
	}
}

func TransformRegisterResponse(u model.User) RegisterResponseBody {
	return RegisterResponseBody{
		Message: "User has been registered successfully",
		ID:      u.ID.String(),
	}
}
