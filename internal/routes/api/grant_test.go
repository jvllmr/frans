package apiRoutes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/testutil"
	"github.com/jvllmr/frans/internal/util"
	"github.com/stretchr/testify/assert"
)

func setupTestGrantRouter(
	testConfig config.Config,
	db *ent.Client,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	group := r.Group("", middlewares...)
	setupGrantGroup(group, testConfig, db)

	return r
}

func createTestGrant(
	t *testing.T,
	db *ent.Client,
	user *ent.User,
	modifier func(q *ent.GrantCreate) *ent.GrantCreate,
) *ent.Grant {
	if modifier == nil {
		modifier = func(q *ent.GrantCreate) *ent.GrantCreate { return q }
	}
	salt := util.GenerateSalt()
	hashedPassword := util.HashPassword("abc123", salt)
	grantValue := modifier(
		db.Grant.Create().
			SetID(uuid.New()).
			SetComment("").
			SetCreatorLang("en").
			SetEmailOnUpload("testmail@vllmr.dev").
			SetExpiryDaysSinceLastUpload(7).
			SetExpiryTotalDays(30).
			SetExpiryTotalUploads(10).
			SetExpiryType("auto").
			SetFileExpiryDaysSinceLastDownload(7).
			SetFileExpiryTotalDays(30).
			SetFileExpiryTotalDownloads(10).
			SetFileExpiryType("auto").
			SetHashedPassword(hashedPassword).SetOwner(user).SetSalt(string(salt)),
	).SaveX(t.Context())

	return grantValue
}

func TestDeleteGrantManually(t *testing.T) {
	cfg := testutil.SetupTestConfig()
	db := testutil.SetupTestDBClient(t)

	testUser := testutil.SetupTestUser(t, db, nil)
	testOwner := testutil.SetupTestUser(t, db, nil)
	testAdmin := testutil.SetupTestAdminUser(t, db, nil)

	rUser := setupTestGrantRouter(cfg, db, testutil.NewTestAuthMiddleware(testUser))
	rOwner := setupTestGrantRouter(cfg, db, testutil.NewTestAuthMiddleware(testOwner))
	rAdmin := setupTestGrantRouter(cfg, db, testutil.NewTestAuthMiddleware(testAdmin))

	testGrant := createTestGrant(t, db, testOwner, nil)
	testGrant2 := createTestGrant(t, db, testOwner, nil)

	reqTestGrant := httptest.NewRequest(http.MethodDelete, "/"+testGrant.ID.String(), nil)

	reqTestGrant2 := httptest.NewRequest(http.MethodDelete, "/"+testGrant2.ID.String(), nil)

	w := httptest.NewRecorder()
	rUser.ServeHTTP(w, reqTestGrant)
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	rOwner.ServeHTTP(w, reqTestGrant)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	rOwner.ServeHTTP(w, reqTestGrant)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, reqTestGrant2)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, reqTestGrant2)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
