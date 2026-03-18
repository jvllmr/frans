package testutil

import (
	"context"
	"testing"

	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/ent"
	"codeberg.org/jvllmr/frans/internal/ent/file"
	"codeberg.org/jvllmr/frans/internal/ent/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupTestUser(
	t *testing.T,
	db *ent.Client,
	modifier func(*ent.UserCreate) *ent.UserCreate,
) *ent.User {
	if modifier == nil {
		modifier = func(uc *ent.UserCreate) *ent.UserCreate { return uc }
	}
	u := modifier(db.User.Create().
		SetID(uuid.New()).
		SetGroups([]string{}).
		SetFullName("Test User").
		SetEmail("testuser@vllmr.dev").
		SetUsername("testuser").
		SetIsAdmin(false)).
		SaveX(t.Context())

	t.Cleanup(func() {
		db.File.Delete().Where(file.HasOwnerWith(user.ID(u.ID))).ExecX(context.Background())
		db.User.DeleteOne(u).ExecX(context.Background())
	})

	return u
}

func SetupTestAdminUser(
	t *testing.T,
	db *ent.Client,
	modifier func(*ent.UserCreate) *ent.UserCreate,
) *ent.User {
	if modifier == nil {
		modifier = func(uc *ent.UserCreate) *ent.UserCreate { return uc }
	}
	u := modifier(db.User.Create().
		SetID(uuid.New()).
		SetGroups([]string{}).
		SetFullName("Test Admin").
		SetEmail("testadmin@vllmr.dev").
		SetUsername("testadmin").
		SetIsAdmin(true)).
		SaveX(t.Context())

	t.Cleanup(func() {
		db.File.Delete().Where(file.HasOwnerWith(user.ID(u.ID))).ExecX(context.Background())
		db.User.DeleteOne(u).ExecX(context.Background())
	})

	return u
}

func NewTestAuthMiddleware(testUser *ent.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.UserGinContext, testUser)
	}
}
