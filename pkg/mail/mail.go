package mail

import (
	"log/slog"

	"github.com/jvllmr/frans/pkg/config"
	gomail "gopkg.in/mail.v2"
)

func sendMail(configValue config.Config, message *gomail.Message) {
	username := ""
	if configValue.SMTPUsername != nil {
		username = *configValue.SMTPUsername
	}
	password := ""
	if configValue.SMTPPassword != nil {
		password = *configValue.SMTPPassword
	}

	dialer := gomail.NewDialer(configValue.SMTPServer, configValue.SMTPPort, username, password)
	message.SetHeader("From", configValue.SMTPFrom)
	if err := dialer.DialAndSend(message); err != nil {
		slog.Error("Could not send mail", "err", err, "recipients", message.GetHeader("To"))
		panic(err)
	} else {
		slog.Info("Sent mail", "recipients", message.GetHeader("To"))
	}
}

func SendFileSharedNotification(configValue config.Config, to string, password *string) {
	// TODO: finish
	message := gomail.NewMessage()

	message.SetHeader("To", to)
	message.SetHeader("Subject", "Download-Ticket")

	message.SetBody("text/plain", "This is the Test Body")
}

func SendFileDownloadNotifiaction(configValue config.Config, to string) {
	// TODO: finish
}
