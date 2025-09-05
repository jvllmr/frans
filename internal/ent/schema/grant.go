package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Grant holds the schema definition for the Grant entity.
type Grant struct {
	ent.Schema
}

// Fields of the Grant.
func (Grant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("comment").Optional().Nillable(),
		field.String("expiry_type"),
		field.String("hashed_password"),
		field.String("salt"),
		field.Time("created_at").
			Default(time.Now),
		field.Uint8("expiry_total_days"),
		field.Uint8("expiry_days_since_last_upload"),
		field.Uint8("expiry_total_uploads"),
		field.String("file_expiry_type"),
		field.Uint8("file_expiry_total_days"),
		field.Uint8("file_expiry_days_since_last_download"),
		field.Uint8("file_expiry_total_downloads"),
		field.Time("last_upload").Nillable().Optional(),
		field.Uint64("times_uploaded").Default(0),
		field.String("email_on_upload").Nillable().Optional(),
		field.String("creator_lang").Default("en"),
	}
}

// Edges of the Grant.
func (Grant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
		edge.From("owner", User.Type).Ref("grants").Unique(),
		edge.To("shareaccesstokens", ShareAccessToken.Type),
	}
}
