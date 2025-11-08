package apiRoutes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func setupTestUserRouter(db *ent.Client, middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	group := r.Group("", middlewares...)
	setupUserGroup(group, db)

	return r
}

func TestFetchMe(t *testing.T) {
	db := testutil.SetupTestDBClient(t)
	testUser := testutil.SetupTestUser(t, db, nil)
	router := setupTestUserRouter(db, testutil.NewTestAuthMiddleware(testUser))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.Bytes()
	var meData services.AdminViewUser
	if err := json.Unmarshal(body, &meData); err != nil {
		log.Fatalf("Could not unmarshal body: %s %v", string(body), err)
	}

	assert.Equal(t, services.ToAdminViewUser(testUser, 0, 0), meData)
}

func TestFetchUsers(t *testing.T) {
	db := testutil.SetupTestDBClient(t)
	users := []*ent.User{
		testutil.SetupTestUser(t, db, nil),
		testutil.SetupTestAdminUser(t, db, nil),
	}
	usersData := make([]services.AdminViewUser, len(users))
	for i, testUser := range users {
		usersData[i] = services.ToAdminViewUser(testUser, 0, 0)
	}
	for _, testUser := range users {
		router := setupTestUserRouter(db, testutil.NewTestAuthMiddleware(testUser))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)
		if testUser.IsAdmin {
			assert.Equal(t, 200, w.Code)
			body := w.Body.Bytes()
			var resultData []services.AdminViewUser
			if err := json.Unmarshal(body, &resultData); err != nil {
				log.Fatalf("Could not unmarshal body: %s %v", string(body), err)
			}

			assert.Equal(t, usersData, resultData)
		} else {
			assert.Equal(t, 403, w.Code)
		}

	}

}
