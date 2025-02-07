package user

type User struct {
	// gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Address string `json:"address"`
}