package sender

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"order-api/configs"
	"order-api/pkg/er"
	"order-api/pkg/event"

	"github.com/jordan-wright/email"
)

type Sender struct {
	EventBus *event.EventBus
	Config    *configs.Config
	Server    string
	Port      string
	Address   string
	TlsConfig *tls.Config
	Auth      smtp.Auth
}

type Addressee struct {
	To string
	Subject string
	Text string
}

func Load(conf *configs.Config, event *event.EventBus) (*Sender, error) {
	// Настройки SMTP
	server := conf.Sender.Address
	port := conf.Sender.Port
	address := server + ":" + port
	err := validate(server, port)
	if err != nil {
		return nil, err
	}
	auth := smtp.PlainAuth("", conf.Sender.Email, conf.Sender.Password, server)

	// // Настроим таймаут подключения
	// dialer := &net.Dialer{
	// 	Timeout:   10 * time.Second,
	// 	KeepAlive: 10 * time.Second,
	// }

	// Настроим TLS
	tlsConfig := &tls.Config{
		ServerName: server,
	}

	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, er.Wrap("Ошибка подключения", err)
	}
	defer conn.Close()

	// Создаем SMTP-клиент
	client, err := smtp.NewClient(conn, server)
	if err != nil {
		return nil, er.Wrap("Ошибка создания клиента", err)
	}
	defer client.Quit()

	// Аутентификация
	if err = client.Auth(auth); err != nil {
		return nil, er.Wrap("Ошибка аутентификации", err)
	}

	return &Sender{
		Config:    conf,
		Server:    server,
		Port:      port,
		Address:   address,
		TlsConfig: tlsConfig,
		Auth:      auth,
		EventBus: event,
	}, nil
}


func validate(server string, port string) error {
	if server == "" {
		return errors.New("server is not specified")
	}
	if port == "" {
		return errors.New("port is not specified")
	}
	return nil
}

func (send *Sender) Listen() {
	for msg := range send.EventBus.Subscribe() {
		if msg.Type == event.EventSendEmail {
			addressee, ok := msg.Data.(Addressee) 
			if !ok {
				log.Fatalln("Bad addressee request ", msg.Data)
				continue
			}

			send.Email(addressee.To, addressee.Subject, addressee.Text)
		}
	}
}

func (send *Sender) Email(to string, subject, text string) error {
	// Настроим письмо
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", send.Config.Sender.Name, send.Config.Sender.Email)
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(text)

	// Отправляем письмо
	if err := e.SendWithTLS(send.Server+":"+send.Port, send.Auth, send.TlsConfig); err != nil {
		return er.Wrap("Ошибка отправки письма", err)
	}
	return nil
}
