package verification

type SendRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}