package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangcollege/sessions"
)

type navigationTestUserRepository struct{}

func (navigationTestUserRepository) CreateUser(string, string, string, string) (int, error) {
	return 0, nil
}

func (navigationTestUserRepository) GetUserByEmail(string) (*User, error) {
	return nil, sql.ErrNoRows
}

func (navigationTestUserRepository) GetUsers() ([]User, error) {
	return nil, nil
}

func (navigationTestUserRepository) Authenticate(string, string) (int, error) {
	return 0, ErrInvalidCredential
}

func TestNavigationRoutes(t *testing.T) {
	session := sessions.New([]byte("test-secret-key-for-navigation"))
	app := &application{
		errorLog:    log.New(io.Discard, "", 0),
		infoLog:     log.New(io.Discard, "", 0),
		userRepo:    navigationTestUserRepository{},
		templateDir: "templates",
		publicPath:  "public",
		session:     session,
	}
	app.tp = NewTemplateRenderer(app.templateDir, false)
	handler := app.routes()

	for _, path := range []string{"/about", "/contact"} {
		t.Run(path, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			handler.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("GET %s returned status %d, want %d", path, response.Code, http.StatusOK)
			}
		})
	}

	t.Run("logout", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/logout", nil)

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusSeeOther {
			t.Fatalf("GET /logout returned status %d, want %d", response.Code, http.StatusSeeOther)
		}
		if location := response.Header().Get("Location"); location != "/" {
			t.Fatalf("GET /logout redirected to %q, want %q", location, "/")
		}
	})
}
