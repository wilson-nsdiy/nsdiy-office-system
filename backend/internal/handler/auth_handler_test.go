//go:build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/domain"
	"oa-nsdiy/backend/internal/service"
	"oa-nsdiy/backend/internal/testutil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthHandler(stub service.UserRepository) *AuthHandler {
	authService := service.NewAuthService(stub, "test-secret-key-for-handler-tests", time.Hour, 7*24*time.Hour)
	return NewAuthHandler(authService)
}

// newTestAuthService creates an AuthService with a stub, pre-hashing a known password.
func newTestAuthService(stub service.UserRepository) *service.AuthService {
	return service.NewAuthService(stub, "test-secret-key", time.Hour, 7*24*time.Hour)
}

func TestLogin_Success_ByUsername(t *testing.T) {
	// Use short salt to keep salt+password under bcrypt 72-byte limit
	testSalt := "s3cret"
	authSvc := newTestAuthService(nil)
	hashedPassword, err := authSvc.HashPassword(context.Background(), "test1234", testSalt)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Nickname:       strPtr("Test"),
				Salt:           testSalt,
				HashedPassword: hashedPassword,
				IsActive:       true,
				TokenVersion:   1,
			},
		},
	}

	auth := setupAuthHandler(stub)

	body := LoginRequest{
		Username: "testuser",
		Password: "test1234",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/login", auth.Login)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				IsActive:       true,
				Salt:           "testsalt",
				HashedPassword: "$2a$12$invalidhashthatwontmatchanypassword",
			},
		},
	}

	auth := setupAuthHandler(stub)

	body := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/login", auth.Login)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestLogin_MissingIdentifier(t *testing.T) {
	auth := setupAuthHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	body := LoginRequest{
		Password: "somepassword",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/login", auth.Login)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLogin_AccountDisabled(t *testing.T) {
	testSalt := "s3cret"
	authSvc := newTestAuthService(nil)
	hashedPassword, _ := authSvc.HashPassword(context.Background(), "test1234", testSalt)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "disableduser",
				Email:          "disabled@example.com",
				Salt:           testSalt,
				HashedPassword: hashedPassword,
				IsActive:       false,
			},
		},
	}

	auth := setupAuthHandler(stub)

	body := LoginRequest{
		Username: "disableduser",
		Password: "test1234",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/login", auth.Login)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	auth := setupAuthHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	body := LoginRequest{
		Username: "nonexistent",
		Password: "password",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/login", auth.Login)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestRequestResetPassword_Success(t *testing.T) {
	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Salt:           "testsalt123456789012345678901234567890",
				HashedPassword: "$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW",
			},
		},
	}

	auth := setupAuthHandler(stub)

	body := ResetPasswordRequest{
		Email: "test@example.com",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/reset/send-code", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/auth/reset/send-code", auth.RequestResetPassword)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}
