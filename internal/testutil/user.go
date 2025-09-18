package testutil

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func SetupTestUser(t *testing.T, db *ent.Client) *ent.User {
	user := db.User.Create().
		SetID(uuid.New()).
		SetGroups([]string{}).
		SetFullName("Test User").
		SetEmail("testuser@vllmr.dev").
		SetUsername("testuser").
		SetIsAdmin(false).
		SaveX(t.Context())

	t.Cleanup(func() {
		db.User.DeleteOne(user).ExecX(context.Background())
	})

	return user
}

func SetupTestAdminUser(t *testing.T, db *ent.Client) *ent.User {
	user := db.User.Create().
		SetID(uuid.New()).
		SetGroups([]string{}).
		SetFullName("Test Admin").
		SetEmail("testadmin@vllmr.dev").
		SetUsername("testadmin").
		SetIsAdmin(true).
		SaveX(t.Context())

	t.Cleanup(func() {
		db.User.DeleteOne(user).ExecX(context.Background())
	})

	return user
}

func NewTestAuthMiddleware(testUser *ent.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.UserGinContext, testUser)
	}
}
