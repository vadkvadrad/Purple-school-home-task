package verify

import (
	"net/smtp"
	"verify-api/internal/user"

	"github.com/jordan-wright/email"
)

type VerifyService struct {
	userRepository *user.UserRepository
}

func NewVerifyService(userRepository *user.UserRepository) *VerifyService {
	return &VerifyService{
		userRepository: userRepository,
	}
}

func(service *VerifyService) GetByEmail(email string) (*user.User, error) {
	return service.userRepository.FindByEmail(email)
}

func(service *VerifyService) GetByHash(hash string) (*user.User, error) {
	return service.userRepository.FindByHash(hash)
}

func(service *VerifyService) Create(email, hash string) (*user.User, error) {
	return service.userRepository.Create(&user.User{
		Email: email,
		Hash: hash,
		IsVerified: false,
	})
}

func(service *VerifyService) Verify(hash string) (*user.User, error) {
	return service.userRepository.Verify(hash)
}

func(service *VerifyService) GenerateHash() (string) {
	hash := user.GenerateHash()
	for {
		existedUser, _ := service.userRepository.FindByHash(hash)
		if existedUser == nil {
			break
		}
		hash = user.GenerateHash()
	}
	return hash
}

func (service *VerifyService) Send(hash string) (error){
	_, err := service.userRepository.FindByHash(hash)
	if err != nil {
		return err
	}
	
	// здесь нужно подставить данные 
	e := email.NewEmail()
	e.From = "Jordan Wright <vadkvadsender@gmail.com>"
	e.To = []string{"mrgrizly2288@gmail.com"}
	e.Subject = "Awesome Subject"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "vadkvadsender@gmail.com", "***", "smtp.gmail.com"))
	return nil
}

func (service *VerifyService) TemSend() error{
	e := email.NewEmail()
	e.From = "Ваше Имя <vadkvadsender@mail.ru>"
	e.To = []string{"vadkvadsender@mail.ru"} // Кому отправляем
	e.Subject = "Тема письма"
	e.Text = []byte("Простое текстовое письмо")
	e.HTML = []byte("<h1>HTML-версия письма</h1>")

	// SMTP-аутентификация
	auth := smtp.PlainAuth("", "vadkvadsender@mail.ru", "***", "smtp.inbox.ru")

	// Отправка письма
	err := e.Send("smtp.inbox.ru:465", auth)
	if err != nil {
		return err
	}
	return nil
}