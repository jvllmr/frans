package mail

import (
	"crypto/tls"
	"log/slog"

	"github.com/wneessen/go-mail"

	"github.com/jvllmr/frans/internal/config"
)

type Mailer struct {
	cfg config.Config
}

func (m *Mailer) sendMail(message *mail.Msg) (err error) {
	defer func() {
		if err != nil {
			slog.Error("Could not send mail", "err", err, "recipients", message.GetToString())
		}
	}()
	username := ""
	if m.cfg.SMTPUsername != nil {
		username = *m.cfg.SMTPUsername
	}
	password := ""
	if m.cfg.SMTPPassword != nil {
		password = *m.cfg.SMTPPassword
	}

	clientOpts := []mail.Option{
		mail.WithPort(m.cfg.SMTPPort),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
	}

	if username != "" {
		clientOpts = append(clientOpts, mail.WithUsername(username))
		clientOpts = append(clientOpts, mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover))
	}

	if password != "" {
		clientOpts = append(clientOpts, mail.WithPassword(password))
	}

	if m.cfg.SMTPInsecureSkipVerify {
		clientOpts = append(clientOpts, mail.WithTLSConfig(
			&tls.Config{
				ServerName:         m.cfg.SMTPServer,
				MinVersion:         mail.DefaultTLSMinVersion,
				InsecureSkipVerify: true,
			},
		))
	}

	dialer, err := mail.NewClient(
		m.cfg.SMTPServer,
		clientOpts...,
	)
	if err != nil {
		slog.Error("Could not setup mail client", "err", err)
		panic(err)
	}
	if err := message.From(m.cfg.SMTPFrom); err != nil {
		return err
	}
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	slog.Info("Sent mail", "recipients", message.GetToString())
	return nil
}

func NewMailer(config config.Config) Mailer {
	return Mailer{cfg: config}
}
