package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	gomail "gopkg.in/mail.v2"
)

func (m *Mailer) SendFileDownloadNotification(
	ctx *gin.Context,
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
	m.sendMail(message)
}
