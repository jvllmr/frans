package services

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

type AdminViewUser struct {
	PublicUser
	ActiveTickets    int   `json:"activeTickets"`
	SubmittedTickets int   `json:"submittedTickets"`
	ActiveGrants     int   `json:"activeGrants"`
	SubmittedGrants  int   `json:"submittedGrants"`
	TotalDataSize    int64 `json:"totalDataSize"`
}

func ToAdminViewUser(user *ent.User, activeTickets int, activeGrants int) AdminViewUser {
	return AdminViewUser{
		PublicUser: PublicUser{
			ID:       user.ID.String(),
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
			Email:    user.Email,
		},
		ActiveTickets:    activeTickets,
		SubmittedTickets: user.SubmittedTickets,
		ActiveGrants:     activeGrants,
		SubmittedGrants:  user.SubmittedGrants,
		TotalDataSize:    user.TotalDataSize,
	}
}
