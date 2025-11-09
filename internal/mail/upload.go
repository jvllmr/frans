package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/wneessen/go-mail"
)

func (m *Mailer) SendFileUploadNotification(
	ctx *gin.Context,

	to string,
	grantValue *ent.Grant,
	files []*ent.File,
) error {
	getTranslation := getTranslationFactory(grantValue.CreatorLang)
	message := mail.NewMsg()

	if err := message.To(to); err != nil {
		return err
	}

	subject := fmt.Sprintf(
		"%s %s",
		getTranslation("subject_upload"),
		grantValue.ID.String(),
	)
	message.Subject(subject)

	bodyTmpl := getTranslation("notification_upload")
	bodyData := map[string]string{
		"GrantName": grantValue.ID.String(),
		"Address":   ctx.ClientIP(),
	}
	var body bytes.Buffer
	if err := template.Must(template.New("").Parse(bodyTmpl)).Execute(&body, bodyData); err != nil {
		panic(err)
	}
	bodyStr := body.String()
	contentsStr := getTranslation("contents") + ":"
	for _, fileValue := range files {
		contentsStr += fmt.Sprintf("\n  - %s", fileValue.Name)
	}
	contentsStr += "\n\n"
	bodyStr = contentsStr + bodyStr

	if grantValue.Comment != nil {
		bodyStr = fmt.Sprintf(
			"%s:\n  %s\n\n",
			getTranslation("grant_comment"),
			*grantValue.Comment,
		) + bodyStr
	}

	message.SetBodyString(mail.TypeTextPlain, bodyStr)
	return m.sendMail(message)
}
