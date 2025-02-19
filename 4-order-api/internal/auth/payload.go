package auth

type LoginRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
	Email string `json:"email" validate:"email"`
}

type LoginResponse struct {
	SessionId string `json:"session_id"`
}

type VerifyRequest struct {
	SessionId string `json:"session_id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type VerifyResponse struct {
	Token string `json:"token"`
}
