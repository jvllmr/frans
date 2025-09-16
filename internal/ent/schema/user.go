package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("username"),
		field.String("full_name"),
		field.String("email"),
		field.Strings("groups"),
		field.Bool("is_admin"),
		field.Time("created_at").
			Default(time.Now),
		field.Int("submitted_tickets").Default(0),
		field.Int("submitted_grants").Default(0),
		field.Int64("totalDataSize").Default(0),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type),
		edge.To("tickets", Ticket.Type),
		edge.To("grants", Grant.Type),
		edge.To("fileinfos", FileData.Type),
	}
}
