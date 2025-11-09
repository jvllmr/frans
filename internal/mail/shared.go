package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	"github.com/wneessen/go-mail"
)

func (m *Mailer) SendTicketSharedNotification(
	ctx *gin.Context,
	ticketService services.TicketService,
	to string,
	locale string,
	ticketValue *ent.Ticket,
	password *string,
) error {
	getTranslation := getTranslationFactory(locale)
	message := mail.NewMsg()

	if err := message.To(to); err != nil {
		return err
	}

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

	message.Subject(subject.String())

	body := fmt.Sprintf(
		"%s: %s",
		getTranslation("url"),
		ticketService.TicketShareLink(ctx, ticketValue),
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

	message.SetBodyString(mail.TypeTextPlain, body)
	return m.sendMail(message)
}

func (m *Mailer) SendGrantSharedNotification(
	ctx *gin.Context,
	grantService services.GrantService,
	to string,
	locale string,
	grantValue *ent.Grant,
	password *string,
) error {
	getTranslation := getTranslationFactory(locale)
	message := mail.NewMsg()

	if err := message.To(to); err != nil {
		return err
	}

	subject := getTranslation("subject_grant")

	message.Subject(subject)

	body := fmt.Sprintf(
		"%s: %s",
		getTranslation("url"),
		grantService.GrantShareLink(ctx, grantValue),
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

	message.SetBodyString(mail.TypeTextPlain, body)
	return m.sendMail(message)
}
