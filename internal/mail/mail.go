package mail

import (
	"log/slog"

	"github.com/jvllmr/frans/internal/config"
	gomail "gopkg.in/mail.v2"
)

type Mailer struct {
	config config.Config
}

func (m *Mailer) sendMail(message *gomail.Message) {
	username := ""
	if m.config.SMTPUsername != nil {
		username = *m.config.SMTPUsername
	}
	password := ""
	if m.config.SMTPPassword != nil {
		password = *m.config.SMTPPassword
	}

	dialer := gomail.NewDialer(m.config.SMTPServer, m.config.SMTPPort, username, password)
	message.SetHeader("From", m.config.SMTPFrom)
	if err := dialer.DialAndSend(message); err != nil {
		slog.Error("Could not send mail", "err", err, "recipients", message.GetHeader("To"))
		panic(err)
	} else {
		slog.Info("Sent mail", "recipients", message.GetHeader("To"))
	}
}

func NewMailer(config config.Config) Mailer {
	return Mailer{config: config}
}
