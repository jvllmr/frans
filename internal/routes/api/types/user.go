package apiTypes

import "github.com/jvllmr/frans/internal/ent"

type PublicUser struct {
	ID       string `json:"id"`
	FullName string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	Email    string `json:"email"`
}

func ToPublicUser(user *ent.User) PublicUser {
	return PublicUser{
		ID:       user.ID.String(),
		FullName: user.FullName,
		IsAdmin:  user.IsAdmin,
		Email:    user.Email,
	}
}
