package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	gomail "gopkg.in/mail.v2"
)

func SendFileUploadNotification(
	ctx *gin.Context,
	configValue config.Config,
	to string,
	grantValue *ent.Grant,
	files []*ent.File,
) {
	getTranslation := getTranslationFactory(grantValue.CreatorLang)
	message := gomail.NewMessage()

	message.SetHeader("To", to)

	subject := fmt.Sprintf(
		"%s %s",
		getTranslation("subject_upload"),
		grantValue.ID.String(),
	)
	message.SetHeader("Subject", subject)

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

	message.SetBody("text/plain", bodyStr)
	sendMail(configValue, message)
}
