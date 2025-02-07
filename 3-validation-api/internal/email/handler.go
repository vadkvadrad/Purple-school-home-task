package email

import "verify-api/configs"

type EmailHandler struct {
	Config *configs.Config
}

type EmailHandlerDeps struct {
	Config *configs.Config
}

// func NewEmailHandler(conf *configs.Config) *EmailHandler {

// }