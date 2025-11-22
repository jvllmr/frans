package apiRoutes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form/v4"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/testutil"
	"github.com/stretchr/testify/assert"
)

var testFormEncoder *form.Encoder = form.NewEncoder()

func setupTestTicketRouter(
	testConfig config.Config,
	db *ent.Client,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	group := r.Group("", middlewares...)
	setupTicketGroup(group, testConfig, db)

	return r
}

func createTestTicket(
	t *testing.T,
	router *gin.Engine,
	inputModifier func(writer *multipart.Writer) int,
) *services.PublicTicket {
	if inputModifier == nil {
		inputModifier = func(writer *multipart.Writer) int { return http.StatusCreated }
	}

	comment := "Test comment"
	email := "test_receiver@vllmr.dev"
	emailOnDownload := "test_creator@vllmr.dev"

	fields := services.TicketFormParams{
		Comment:                     &comment,
		Email:                       &email,
		Password:                    "abc123",
		EmailPassword:               true,
		ExpiryType:                  "auto",
		ExpiryTotalDays:             30,
		ExpiryDaysSinceLastDownload: 7,
		ExpiryTotalDownloads:        10,
		EmailOnDownload:             &emailOnDownload,
		CreatorLang:                 "en",
		ReceiverLang:                "en",
	}

	encodedFields, err := testFormEncoder.Encode(&fields)
	if err != nil {
		log.Fatalf("encode form: %v", err)
	}

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for k, values := range encodedFields {
		for _, v := range values {
			err = writer.WriteField(k, v)
			if err != nil {
				log.Fatalf("write form field: %v", err)
			}
		}
	}

	expectedStatus := inputModifier(writer)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
	if w.Code != http.StatusCreated {
		return nil
	}
	var publicTicket services.PublicTicket
	if err = json.Unmarshal(w.Body.Bytes(), &publicTicket); err != nil {
		log.Fatalf("unmarshal ticket: %v", err)
	}
	return &publicTicket

}

func TestCreateBasicTicket(t *testing.T) {
	db := testutil.SetupTestDBClient(t)
	testUser := testutil.SetupTestUser(t, db, nil)
	testConfig := testutil.SetupTestConfig()
	router := setupTestTicketRouter(testConfig, db, testutil.NewTestAuthMiddleware(testUser))

	ticketInputModifier := func(writer *multipart.Writer) int {
		partWriter, _ := writer.CreateFormFile("files[]", "test.txt")
		io.Copy(partWriter, strings.NewReader("This is a test file. Say hello!"))

		return http.StatusCreated
	}

	newTicket := createTestTicket(t, router, ticketInputModifier)

	assert.Equal(t, 1, len(newTicket.Files))
}

func TestFetchTickets(t *testing.T) {
	db := testutil.SetupTestDBClient(t)
	testUser := testutil.SetupTestUser(t, db, nil)
	testAdminUser := testutil.SetupTestAdminUser(t, db, nil)
	testConfig := testutil.SetupTestConfig()
	router := setupTestTicketRouter(testConfig, db, testutil.NewTestAuthMiddleware(testUser))
	adminRouter := setupTestTicketRouter(
		testConfig,
		db,
		testutil.NewTestAuthMiddleware(testAdminUser),
	)

	ticketInputModifier := func(writer *multipart.Writer) int {
		partWriter, _ := writer.CreateFormFile("files[]", "test.txt")
		io.Copy(partWriter, strings.NewReader("This is a test file. Say hello!"))

		return http.StatusCreated
	}

	newTicket := createTestTicket(t, router, ticketInputModifier)
	assert.Equal(t, 1, len(newTicket.Files))

	newAdminTicket := createTestTicket(t, adminRouter, ticketInputModifier)
	assert.Equal(t, 1, len(newAdminTicket.Files))
	assert.NotEqual(t, newTicket.ID, newAdminTicket.ID)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resultTickets []*services.PublicTicket
	if err := json.Unmarshal(w.Body.Bytes(), &resultTickets); err != nil {
		log.Fatalf("unmsarshal tickets: %v", err)
	}
	assert.Equal(t, 1, len(resultTickets))
	assert.Equal(t, newTicket, resultTickets[0])

	wAdmin := httptest.NewRecorder()
	adminRouter.ServeHTTP(wAdmin, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var adminTickets []*services.PublicTicket
	if err := json.Unmarshal(wAdmin.Body.Bytes(), &adminTickets); err != nil {
		log.Fatalf("unmsarshal tickets: %v", err)
	}

	assert.Equal(t, 2, len(adminTickets))
	assert.NotEqual(t, len(resultTickets), len(adminTickets), wAdmin.Body.String())

}

func TestDeleteTicketManually(t *testing.T) {
	cfg := testutil.SetupTestConfig()
	db := testutil.SetupTestDBClient(t)

	testUser := testutil.SetupTestUser(t, db, nil)
	testOwner := testutil.SetupTestUser(t, db, nil)
	testAdmin := testutil.SetupTestAdminUser(t, db, nil)

	rUser := setupTestTicketRouter(cfg, db, testutil.NewTestAuthMiddleware(testUser))
	rOwner := setupTestTicketRouter(cfg, db, testutil.NewTestAuthMiddleware(testOwner))
	rAdmin := setupTestTicketRouter(cfg, db, testutil.NewTestAuthMiddleware(testAdmin))

	testTicket := createTestTicket(t, rOwner, nil)
	testTicket2 := createTestTicket(t, rOwner, nil)

	reqTestTicket := httptest.NewRequest(http.MethodDelete, "/"+testTicket.ID.String(), nil)

	reqTestTicket2 := httptest.NewRequest(http.MethodDelete, "/"+testTicket2.ID.String(), nil)

	w := httptest.NewRecorder()
	rUser.ServeHTTP(w, reqTestTicket)
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	rOwner.ServeHTTP(w, reqTestTicket)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	rOwner.ServeHTTP(w, reqTestTicket)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, reqTestTicket2)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, reqTestTicket2)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
