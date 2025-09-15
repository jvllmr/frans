package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("name"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("last_download").Nillable().Optional(),
		field.Uint64("times_downloaded").Default(0),
		field.String("expiry_type"),
		field.Uint8("expiry_total_days"),
		field.Uint8("expiry_days_since_last_download"),
		field.Uint8("expiry_total_downloads"),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tickets", Ticket.Type).Ref("files"),
		edge.From("grants", Grant.Type).Ref("files"),
		edge.To("data", FileData.Type).Unique().Required(),
	}
}
