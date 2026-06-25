//go:build unit

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/domain"
	"oa-nsdiy/backend/internal/service"
	"oa-nsdiy/backend/internal/setup"
	"oa-nsdiy/backend/internal/testutil"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})
	return dir
}

func newSetupHandler(stub service.UserRepository) *SetupHandler {
	authService := service.NewAuthService(stub, "test-secret", time.Hour, 7*24*time.Hour)
	return NewSetupHandler(authService)
}

func TestNeedsSetupDoubleCheck_FileExists_NoDBCheck(t *testing.T) {
	setupTestDir(t)
	if err := setup.CreateInstallLock(); err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{1: testutil.NewTestUser()},
	})

	needsSetup, err := h.needsSetupDoubleCheck(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if needsSetup {
		t.Error("expected needsSetup=false when .installed file exists")
	}
}

func TestNeedsSetupDoubleCheck_NoFile_NoUsers(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	needsSetup, err := h.needsSetupDoubleCheck(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !needsSetup {
		t.Error("expected needsSetup=true when no file and no users")
	}
}

func TestNeedsSetupDoubleCheck_NoFile_HasUsers(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{1: testutil.NewTestUser()},
	})

	needsSetup, err := h.needsSetupDoubleCheck(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if needsSetup {
		t.Error("expected needsSetup=false when users exist in DB")
	}
}

func TestNeedsSetupDoubleCheck_NoFile_DatabaseError(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
		Err:   fmt.Errorf("database connection failed"),
	})

	_, err := h.needsSetupDoubleCheck(t.Context())
	if err == nil {
		t.Error("expected error when database query fails")
	}
}

func TestStatusEndpoint_NeedsSetup(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	req := httptest.NewRequest("GET", "/api/setup/status", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/setup/status", h.Status)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Data struct {
			NeedsSetup bool `json:"needsSetup"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Data.NeedsSetup {
		t.Error("expected needsSetup=true")
	}
}

func TestStatusEndpoint_AlreadyInstalled(t *testing.T) {
	setupTestDir(t)
	if err := setup.CreateInstallLock(); err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	req := httptest.NewRequest("GET", "/api/setup/status", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/setup/status", h.Status)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Data struct {
			NeedsSetup bool `json:"needsSetup"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Data.NeedsSetup {
		t.Error("expected needsSetup=false when lock file exists")
	}
}

func TestSetupGuard_BlocksWhenInstalled(t *testing.T) {
	setupTestDir(t)
	if err := setup.CreateInstallLock(); err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	req := httptest.NewRequest("POST", "/api/setup/admin", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(h.SetupGuard())
	router.POST("/api/setup/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestSetupGuard_AllowsWhenNotInstalled(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	req := httptest.NewRequest("POST", "/api/setup/admin", nil)
	w := httptest.NewRecorder()

	guardHit := false
	router := gin.New()
	router.Use(h.SetupGuard())
	router.POST("/api/setup/admin", func(c *gin.Context) {
		guardHit = true
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.ServeHTTP(w, req)

	if !guardHit {
		t.Error("expected handler to be called when not installed")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCreateAdmin_DoubleCheck_PreventsReinstall(t *testing.T) {
	setupTestDir(t)
	if err := setup.CreateInstallLock(); err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	})

	body := CreateAdminRequest{
		Username: "admin",
		Email:    "admin@test.com",
		Password: "password123",
		Nickname: "Admin",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/setup/admin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/setup/admin", h.CreateAdmin)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateAdmin_DoubleCheck_PreventsReinstallByDB(t *testing.T) {
	setupTestDir(t)

	h := newSetupHandler(&testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{1: testutil.NewTestUser()},
	})

	body := CreateAdminRequest{
		Username: "admin2",
		Email:    "admin2@test.com",
		Password: "password123",
		Nickname: "Admin2",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/setup/admin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/setup/admin", h.CreateAdmin)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 (users exist in DB), got %d: %s", w.Code, w.Body.String())
	}
}
