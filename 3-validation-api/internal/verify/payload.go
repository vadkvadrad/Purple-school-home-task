package verify

type SendRequest struct {
	Email    string `json:"email" validate:"required,email"`
}
