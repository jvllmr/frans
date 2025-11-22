package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/wneessen/go-mail"
)

func (m *Mailer) sendDeletionNotification(
	to string,
	lang string,
	entity string,
	baseUrl string,
) error {
	getTranslation := getTranslationFactory(lang)
	message := mail.NewMsg()

	if err := message.To(to); err != nil {
		return err
	}

	subject := fmt.Sprintf(
		"%s %s",
		getTranslation("subject_deletion"),
		entity,
	)
	message.Subject(subject)

	bodyTmpl := getTranslation("notification_deletion")
	bodyData := map[string]string{
		"Entity":  entity,
		"BaseUrl": baseUrl,
	}
	var body bytes.Buffer
	if err := template.Must(template.New("").Parse(bodyTmpl)).Execute(&body, bodyData); err != nil {
		panic(err)
	}
	message.SetBodyString(mail.TypeTextPlain, body.String())
	return m.sendMail(message)
}

func (m *Mailer) SendTicketDeletionNotification(t *ent.Ticket, baseUrl string) error {
	getTranslation := getTranslationFactory(t.CreatorLang)

	entityTmpl := getTranslation("entity_ticket")
	entityData := map[string]string{
		"ID": t.ID.String(),
	}
	var entity bytes.Buffer
	if err := template.Must(template.New("").Parse(entityTmpl)).Execute(&entity, entityData); err != nil {
		panic(err)
	}

	return m.sendDeletionNotification(t.Edges.Owner.Email, t.CreatorLang, entity.String(), baseUrl)
}

func (m *Mailer) SendGrantDeletionNotification(g *ent.Grant, baseUrl string) error {
	getTranslation := getTranslationFactory(g.CreatorLang)

	entityTmpl := getTranslation("entity_grant")
	entityData := map[string]string{
		"ID": g.ID.String(),
	}
	var entity bytes.Buffer
	if err := template.Must(template.New("").Parse(entityTmpl)).Execute(&entity, entityData); err != nil {
		panic(err)
	}

	return m.sendDeletionNotification(g.Edges.Owner.Email, g.CreatorLang, entity.String(), baseUrl)
}

func (m *Mailer) SendFileDeletionNotification(f *ent.File, baseUrl string) error {
	lang := "en"

	if f.Edges.Grant != nil {
		lang = f.Edges.Grant.CreatorLang
	}

	if f.Edges.Ticket != nil {
		lang = f.Edges.Ticket.CreatorLang
	}

	getTranslation := getTranslationFactory(lang)

	entityTmpl := getTranslation("entity_file")
	entityData := map[string]string{
		"ID": f.ID.String(),
	}
	var entity bytes.Buffer
	if err := template.Must(template.New("").Parse(entityTmpl)).Execute(&entity, entityData); err != nil {
		panic(err)
	}

	return m.sendDeletionNotification(f.Edges.Owner.Email, lang, entity.String(), baseUrl)
}
