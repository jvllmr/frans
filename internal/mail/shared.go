package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/util"
	gomail "gopkg.in/mail.v2"
)

func SendTicketSharedNotification(
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
		util.TicketShareLink(ctx, configValue, ticketValue),
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
		body = fmt.Sprintf(
			"%s:\n  %s\n\n",
			getTranslation("ticket_comment"),
			*ticketValue.Comment,
		) + body
	}

	message.SetBody("text/plain", body)
	sendMail(configValue, message)
}

func SendGrantSharedNotification(
	ctx *gin.Context,
	configValue config.Config,
	to string,
	locale string,
	grantValue *ent.Grant,
	password *string,
) {
	getTranslation := getTranslationFactory(locale)
	message := gomail.NewMessage()

	message.SetHeader("To", to)

	subject := getTranslation("subject_grant")

	message.SetHeader("Subject", subject)

	body := fmt.Sprintf(
		"%s: %s",
		getTranslation("url"),
		util.GrantShareLink(ctx, configValue, grantValue),
	)

	if password != nil {
		body += fmt.Sprintf("\n%s: %s", getTranslation("password"), *password)
	}

	if grantValue.Comment != nil {
		body = fmt.Sprintf(
			"%s:\n  %s\n\n",
			getTranslation("grant_comment"),
			*grantValue.Comment,
		) + body
	}

	message.SetBody("text/plain", body)
	sendMail(configValue, message)
}
