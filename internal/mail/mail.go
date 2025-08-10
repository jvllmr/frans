package mail

import (
	"bytes"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
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

func SendFileSharedNotification(
	ctx *gin.Context,
	configValue config.Config,
	to string,
	locale string,
	ticketValue *ent.Ticket,
	password *string,
) {
	getTranslation := getTranslationFactory(locale)
	message := gomail.NewMessage()

	message.SetHeader("To", to)

	subjectTmpl := getTranslation("subject_ticket")

	subjectFileName := ticketValue.ID.String()

	if len(ticketValue.Edges.Files) == 1 {
		subjectFileName = ticketValue.Edges.Files[0].Name
	}

	subjectData := map[string]string{
		"FileName": subjectFileName,
	}
	var subject bytes.Buffer
	if err := template.Must(template.New("").Parse(subjectTmpl)).Execute(&subject, subjectData); err != nil {
		panic(err)
	}

	message.SetHeader("Subject", subject.String())

	body := fmt.Sprintf(
		"%s: %s",
		getTranslation("url"),
		util.GetTicketShareLink(ctx, configValue, ticketValue),
	)

	if password != nil {
		body += fmt.Sprintf("\n%s: %s", getTranslation("password"), *password)
	}

	if len(ticketValue.Edges.Files) > 1 {
		contentsStr := getTranslation("contents") + ":"
		for _, file := range ticketValue.Edges.Files {
			contentsStr += fmt.Sprintf("\n  - %s", file.Name)
		}
		contentsStr += "\n\n"
		body = contentsStr + body
	}

	if ticketValue.Comment != nil {
		body = fmt.Sprintf("%s:\n  %s\n\n", getTranslation("comment"), *ticketValue.Comment) + body
	}

	message.SetBody("text/plain", body)
	sendMail(configValue, message)
}

func SendFileDownloadNotification(
	ctx *gin.Context,
	configValue config.Config,
	to string,
	ticketValue *ent.Ticket,
	fileValue *ent.File,
) {
	getTranslation := getTranslationFactory(ticketValue.CreatorLang)
	message := gomail.NewMessage()

	message.SetHeader("To", to)

	subject := fmt.Sprintf(
		"%s %s (%s)",
		getTranslation("subject_download"),
		ticketValue.ID.String(),
		fileValue.Name,
	)
	message.SetHeader("Subject", subject)

	bodyTmpl := getTranslation("notification_download")
	bodyData := map[string]string{
		"FileName":   fileValue.Name,
		"TicketName": ticketValue.ID.String(),
		"Address":    ctx.ClientIP(),
	}
	var body bytes.Buffer
	if err := template.Must(template.New("").Parse(bodyTmpl)).Execute(&body, bodyData); err != nil {
		panic(err)
	}
	message.SetBody("text/plain", body.String())
	sendMail(configValue, message)
}
