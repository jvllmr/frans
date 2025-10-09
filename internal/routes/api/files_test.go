package apiRoutes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func setupTestFileRouter(
	testConfig config.Config,
	db *ent.Client,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	group := r.Group("", middlewares...)
	setupFileGroup(group, testConfig, db)

	return r
}

func TestFetchFile(t *testing.T) {
	cfg := testutil.SetupTestConfig()
	db := testutil.SetupTestDBClient(t)

	testUser := testutil.SetupTestUser(t, db, nil)
	testFile := testutil.SetupTestFile(
		t,
		cfg,
		db,
		"test.txt",
		"Hello there!",
		testUser,
		"single",
		0,
		0,
		1,
	)

	r := setupTestFileRouter(cfg, db, testutil.NewTestAuthMiddleware(testUser))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testFile.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFetchReceivedFiles(t *testing.T) {
	cfg := testutil.SetupTestConfig()
	db := testutil.SetupTestDBClient(t)

	testUser := testutil.SetupTestUser(t, db, nil)
	testAdmin := testutil.SetupTestAdminUser(t, db, nil)

	userGrantValue := createTestGrant(t, db, testUser, nil)
	adminGrantValue := createTestGrant(t, db, testAdmin, nil)

	testFile := testutil.SetupTestFile(
		t,
		cfg,
		db,
		"test.txt",
		"Hello there!",
		testUser,
		"single",
		0,
		0,
		1,
	)
	db.File.UpdateOne(testFile).AddGrants(userGrantValue).SaveX(t.Context())

	testFileAdmin := testutil.SetupTestFile(
		t,
		cfg,
		db,
		"test2.txt",
		"Hello there!",
		testAdmin,
		"single",
		0,
		0,
		1,
	)
	db.File.UpdateOne(testFileAdmin).AddGrants(adminGrantValue).SaveX(t.Context())

	rUser := setupTestFileRouter(cfg, db, testutil.NewTestAuthMiddleware(testUser))
	rAdmin := setupTestFileRouter(cfg, db, testutil.NewTestAuthMiddleware(testAdmin))

	req := httptest.NewRequest(http.MethodGet, "/received", nil)
	w := httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resultFiles []*services.PublicFile
	if err := json.Unmarshal(w.Body.Bytes(), &resultFiles); err != nil {
		log.Fatalf("unmsarshal files: %v", err)
	}
	assert.Equal(t, 1, len(resultFiles))
	assert.Equal(t, testFile.ID, resultFiles[0].Id)

	wAdmin := httptest.NewRecorder()
	rAdmin.ServeHTTP(wAdmin, req)
	assert.Equal(t, http.StatusOK, wAdmin.Code)
	var resultFilesAdmin []*services.PublicFile
	if err := json.Unmarshal(wAdmin.Body.Bytes(), &resultFilesAdmin); err != nil {
		log.Fatalf("unmsarshal files: %v", err)
	}
	assert.Equal(t, 1, len(resultFilesAdmin))
	assert.Equal(t, testFileAdmin.ID, resultFilesAdmin[0].Id)
}
