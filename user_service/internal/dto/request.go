package dto

type RegisterRequestBody struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyRequestBody struct {
	Token string `json:"token"`
}
