package mailserv

import (
	"bytes"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/repository/mailrepository"
	_ "embed"
	"log"
	"text/template"

	"github.com/go-mail/mail/v2"
	"gorm.io/gorm"
)

var (
	//go:embed template.html
	msgTemplate      string
	ErrSendedPrevius = "El mensaje se envio previamete"
)

type MailService interface {
	Send(to, subject, body string) error
	HandleMsg(purchase *model.Purchase)
	FormatMsg(peticion *model.EmailMsg) ([]byte, error)
}

type MailServ struct {
	t      *template.Template
	repo   repository.MailRepository
	mailer Mailer
	sender string
}

func New(db *gorm.DB, config *MailerConfig) MailService {
	dailer := mail.NewDialer(
		config.Host,
		config.Port,
		config.Username,
		config.Password,
	)

	dailer.Timeout = config.Timeout

	mailer := Mailer{
		dailer: dailer,
		config: *config,
	}

	return &MailServ{
		t: template.Must(template.New("template").
			Parse(msgTemplate)),
		repo:   mailrepository.New(db),
		sender: config.Sender,
		mailer: mailer,
	}
}

func (ms *MailServ) HandleMsg(purchase *model.Purchase) {
	emailMsgStruct, err := ms.repo.GetMsg(purchase)
	if err != nil {
		if err.Error() != ErrSendedPrevius {
			log.Println(err.Error())
			return
		}
	}
	templat, err := ms.FormatMsg(emailMsgStruct)
	if err != nil {
		log.Println("Error al formatear el mensaje de correo")
		return
	}

	ms.Send(emailMsgStruct.To, "Farmu Compra", string(templat))
}

func (ms *MailServ) FormatMsg(peticion *model.EmailMsg) ([]byte, error) {
	var b bytes.Buffer
	if err := ms.t.
		Execute(&b, peticion); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (ms *MailServ) Send(to, subject, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("From", ms.sender)
	msg.SetBody("text/html", body)
	return ms.mailer.dailer.DialAndSend(msg)
}
