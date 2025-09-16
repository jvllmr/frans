package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FileData holds the schema definition for the FileData entity.
type FileData struct {
	ent.Schema
}

// Fields of the FileData.
func (FileData) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Uint64("size"),
	}
}

// Edges of the FileData.
func (FileData) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("fileinfos").Required(),
		edge.From("files", File.Type).Ref("data"),
	}
}
