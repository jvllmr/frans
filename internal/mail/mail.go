package mail

import (
	"log/slog"

	"github.com/wneessen/go-mail"

	"github.com/jvllmr/frans/internal/config"
)

type Mailer struct {
	config config.Config
}

func (m *Mailer) sendMail(message *mail.Msg) (err error) {
	defer func() {
		if err != nil {
			slog.Error("Could not send mail", "err", err, "recipients", message.GetToString())
		}
	}()
	username := ""
	if m.config.SMTPUsername != nil {
		username = *m.config.SMTPUsername
	}
	password := ""
	if m.config.SMTPPassword != nil {
		password = *m.config.SMTPPassword
	}

	dialer, err := mail.NewClient(
		m.config.SMTPServer,
		mail.WithPort(m.config.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
	)
	if err != nil {
		slog.Error("Could not setup mail client", "err", err)
		panic(err)
	}
	if err := message.From(m.config.SMTPFrom); err != nil {
		return err
	}
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	slog.Info("Sent mail", "recipients", message.GetToString())
	return nil
}

func NewMailer(config config.Config) Mailer {
	return Mailer{config: config}
}
