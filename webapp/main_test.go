package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/stretchr/testify/assert"
)

func createTestApp() *App { // helper
	_ = os.Setenv("AUTH0_DOMAIN", "example.auth0.com")
	_ = os.Setenv("AUTH0_CLIENT_ID", "client_id")
	_ = os.Setenv("AUTH0_CLIENT_SECRET", "client_secret")

	app := NewDemoApp()
	store := cookie.NewStore([]byte("secret"))
	app.router.Use(sessions.Sessions("auth-session", store))

	return app
}

func TestLoginHandler(t *testing.T) { // todo
	app := createTestApp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("email=test@example.com&password=password123")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateProfileHandler(t *testing.T) {
	app := createTestApp()
	w := httptest.NewRecorder()
	profileData := profile{
		Age:                      30,
		Name:                     "Test User",
		Email:                    "test@example.com",
		FavouriteColor:           "blue",
		FavouriteOperatingSystem: "linux",
	}
	body, _ := json.Marshal(profileData)
	req, _ := http.NewRequest("POST", "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteProfileHandler(t *testing.T) {
	app := createTestApp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/profile", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLogoutHandler(t *testing.T) {
	app := createTestApp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/logout", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestIsAuthenticated(t *testing.T) {
	app := createTestApp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestAuth0Initialization(t *testing.T) {
	app := createTestApp()
	assert.NotNil(t, app.auth)
}

func TestMemDBInitialization(t *testing.T) {
	app := createTestApp()
	assert.NotNil(t, app.db)
}

func TestInsertProfileInMemDB(t *testing.T) {
	app := createTestApp()
	txn := app.db.Txn(true)
	profile := profile{
		Age:                      25,
		Name:                     "Test User",
		Email:                    "test@example.com",
		FavouriteColor:           "red",
		FavouriteOperatingSystem: "linux",
	}
	err := txn.Insert(PROFILE, profile)
	assert.Nil(t, err)
	txn.Commit()
}

func TestRetrieveProfileFromMemDB(t *testing.T) {
	app := createTestApp()
	txn := app.db.Txn(true)
	pro := profile{
		Age:                      25,
		Name:                     "Test User",
		Email:                    "test@example.com",
		FavouriteColor:           "red",
		FavouriteOperatingSystem: "linux",
	}
	txn.Insert(PROFILE, pro)
	txn.Commit()

	txn = app.db.Txn(false)
	raw, err := txn.First(PROFILE, ID, "test@example.com")
	assert.Nil(t, err)
	assert.NotNil(t, raw)
	retrievedProfile := raw.(*profile)
	assert.Equal(t, "Test User", retrievedProfile.Name)
}

func TestDeleteProfileFromMemDB(t *testing.T) {
	app := createTestApp()
	txn := app.db.Txn(true)
	pro := profile{
		Age:                      25,
		Name:                     "Test User",
		Email:                    "test@example.com",
		FavouriteColor:           "red",
		FavouriteOperatingSystem: "linux",
	}
	txn.Insert(PROFILE, pro)
	txn.Commit()

	txn = app.db.Txn(true)
	err := txn.Delete(PROFILE, pro)
	assert.Nil(t, err)
	txn.Commit()

	txn = app.db.Txn(false)
	raw, err := txn.First(PROFILE, ID, "test@example.com")
	assert.Nil(t, err)
	assert.Nil(t, raw)
}

func TestAppInitialization(t *testing.T) {
	app := createTestApp()
	assert.NotNil(t, app.router)
	assert.NotNil(t, app.db)
	assert.NotNil(t, app.auth)
	assert.NotNil(t, app.validate)
}
