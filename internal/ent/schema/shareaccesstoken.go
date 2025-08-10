package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ShareAccessToken holds the schema definition for the ShareAccessToken entity.
type ShareAccessToken struct {
	ent.Schema
}

// Fields of the ShareAccessToken.
func (ShareAccessToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Time("expiry"),
	}
}

// Edges of the ShareAccessToken.
func (ShareAccessToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ticket", Ticket.Type).Ref("shareaccesstokens").Unique(),
	}
}
