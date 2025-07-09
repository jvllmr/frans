package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Ticket holds the schema definition for the Ticket entity.
type Ticket struct {
	ent.Schema
}

// Fields of the Ticket.
func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("comment").Optional().Nillable(),
		field.String("expiryType"),
		field.String("hashed_password"),
		field.String("salt"),
		field.Time("created_at").
			Default(time.Now),
		field.Uint8("expiry_total_days"),
		field.Uint8("expiry_days_since_last_download"),
		field.Uint8("expiry_total_downloads"),
		field.String("email_on_download").Nillable().Optional(),
	}
}

// Edges of the Ticket.
func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
		edge.From("owner", User.Type).Ref("tickets").Unique(),
	}
}
