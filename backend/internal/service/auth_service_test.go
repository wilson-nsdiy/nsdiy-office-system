//go:build unit
// +build unit

package service

import (
	"context"
	"testing"
	"time"

	"oa-nsdiy/backend/internal/domain"
	"oa-nsdiy/backend/internal/testutil"
)

// shortSalt is a manually constructed salt that keeps salt+password under bcrypt's 72-byte limit.
const shortSalt = "s3cret"

// testPasswordHash uses a real service to hash with a short salt.
func testPasswordHash(auth *AuthService, password string) string {
	hash, err := auth.HashPassword(context.Background(), password, shortSalt)
	if err != nil {
		panic("test setup failed: " + err.Error())
	}
	return hash
}

func TestAuthService_ChangePassword_Incorrect(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hashedPassword := testPasswordHash(auth, "old-password")

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Salt:           shortSalt,
				HashedPassword: hashedPassword,
				TokenVersion:   1,
			},
		},
	}

	auth.userRepo = stub

	// Use exactly 8-char new password: salt(64) + pass(8) = 72 = bcrypt limit
	err := auth.ChangePassword(context.Background(), 1, "wrong-password", "Passw0rd")
	if err == nil {
		t.Fatal("expected error for incorrect old password, got nil")
	}
	if err.Error() != "old password is incorrect" {
		t.Errorf("expected 'old password is incorrect', got '%v'", err)
	}
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hashedPassword := testPasswordHash(auth, "old-password")

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Salt:           shortSalt,
				HashedPassword: hashedPassword,
				TokenVersion:   1,
			},
		},
	}

	auth.userRepo = stub

	// Use exactly 8-char new password: salt(64) + pass(8) = 72 = bcrypt limit
	err := auth.ChangePassword(context.Background(), 1, "old-password", "Passw0rd")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	user := stub.Users[1]
	if user.Salt == shortSalt {
		t.Error("expected salt to change after password update")
	}
	if user.HashedPassword == hashedPassword {
		t.Error("expected hashedPassword to change after password update")
	}

	if !auth.VerifyPassword(context.Background(), "Passw0rd", user.HashedPassword, user.Salt) {
		t.Error("new password should pass verification")
	}

	if user.TokenVersion != 2 {
		t.Errorf("expected TokenVersion=2, got %d", user.TokenVersion)
	}
}

func TestAuthService_ValidatePasswordStrength(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)
	ctx := context.Background()

	tests := []struct {
		password string
		valid    bool
	}{
		{"short", false},
		{"12345678", false},
		{"abcdefgh", false},
		{"abcd1234", true},
		{"Password123", true},
		{"A1b2C3d4", true},
		{"abc1", false},
		{"abcdefgh12345678", true},
	}

	for _, tt := range tests {
		result := auth.ValidatePasswordStrength(ctx, tt.password)
		if result != tt.valid {
			t.Errorf("ValidatePasswordStrength(%q) = %v, want %v", tt.password, result, tt.valid)
		}
	}
}

func TestAuthService_ConfirmResetPassword_ExpiredCode(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hashedPassword := testPasswordHash(auth, "old-password")
	expiredTime := time.Now().Add(-1 * time.Hour)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:                        1,
				Username:                  "testuser",
				Email:                     "test@example.com",
				Salt:                      shortSalt,
				HashedPassword:            hashedPassword,
				TokenVersion:              1,
				VerificationCode:          strPtr("888888"),
				VerificationCodeExpiresAt: &expiredTime,
			},
		},
	}

	auth.userRepo = stub

	err := auth.ConfirmResetPassword(context.Background(), "test@example.com", "888888", "NewPass123")
	if err == nil {
		t.Fatal("expected error for expired code, got nil")
	}
	if err.Error() != "verification code has expired" {
		t.Errorf("expected 'verification code has expired', got '%v'", err)
	}
}

func TestAuthService_ConfirmResetPassword_InvalidCode(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hashedPassword := testPasswordHash(auth, "old-password")
	validTime := time.Now().Add(30 * time.Minute)
	code := "123456"

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:                        1,
				Username:                  "testuser",
				Email:                     "test@example.com",
				Salt:                      shortSalt,
				HashedPassword:            hashedPassword,
				TokenVersion:              1,
				VerificationCode:          &code,
				VerificationCodeExpiresAt: &validTime,
			},
		},
	}

	auth.userRepo = stub

	err := auth.ConfirmResetPassword(context.Background(), "test@example.com", "000000", "NewPass123")
	if err == nil {
		t.Fatal("expected error for invalid code, got nil")
	}
	if err.Error() != "invalid verification code" {
		t.Errorf("expected 'invalid verification code', got '%v'", err)
	}
}

func TestAuthService_RequestResetPassword_Success(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hashedPassword := testPasswordHash(auth, "password")

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Salt:           shortSalt,
				HashedPassword: hashedPassword,
			},
		},
	}

	auth.userRepo = stub

	code, expiresAt, err := auth.RequestResetPassword(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6-digit code, got %q", code)
	}

	if expiresAt.Before(time.Now()) {
		t.Error("expected valid expiry time")
	}

	user := stub.Users[1]
	if user.VerificationCode == nil || *user.VerificationCode != code {
		t.Error("verification code was not stored on user")
	}
}

func TestAuthService_RequestResetPassword_UserNotFound(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{},
	}

	auth.userRepo = stub

	_, _, err := auth.RequestResetPassword(context.Background(), "nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for unknown email, got nil")
	}
}

func TestAuthService_GenerateSalt(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	salt1, err := auth.GenerateSalt(context.Background())
	if err != nil {
		t.Fatalf("failed to generate salt: %v", err)
	}

	salt2, err := auth.GenerateSalt(context.Background())
	if err != nil {
		t.Fatalf("failed to generate salt: %v", err)
	}

	if salt1 == salt2 {
		t.Error("two generated salts should be different")
	}
	if len(salt1) != 64 {
		t.Errorf("expected 64-char hex salt, got %d chars", len(salt1))
	}
}

func TestAuthService_HashToken(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	hash1 := auth.HashToken(context.Background(), "test-token")
	hash2 := auth.HashToken(context.Background(), "test-token")
	hash3 := auth.HashToken(context.Background(), "different-token")

	if hash1 != hash2 {
		t.Error("same input should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("different inputs should produce different hashes")
	}
	if len(hash1) != 64 {
		t.Errorf("expected 64-char hex hash (SHA256), got %d chars", len(hash1))
	}
}

func TestAuthService_GenerateApiToken(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	token, prefix, hash := auth.GenerateApiToken(context.Background())

	if len(token) != 40 {
		t.Errorf("expected 40-char hex token (20 bytes), got %d chars", len(token))
	}
	if len(prefix) != 8 {
		t.Errorf("expected 8-char prefix, got %q (%d chars)", prefix, len(prefix))
	}
	if len(hash) != 64 {
		t.Errorf("expected 64-char hash (SHA256), got %d chars", len(hash))
	}
	if token[:8] != prefix {
		t.Error("token prefix should match the first 8 chars")
	}
}

func TestAuthService_CreateAndVerifyToken(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	accessToken, err := auth.CreateAccessToken(context.Background(), 42, "testuser", 1)
	if err != nil {
		t.Fatalf("failed to create access token: %v", err)
	}

	claims, err := auth.VerifyToken(context.Background(), accessToken)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("expected UserID=42, got %d", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("expected Username=testuser, got %s", claims.Username)
	}
	if claims.TokenType != "access" {
		t.Errorf("expected TokenType=access, got %s", claims.TokenType)
	}
	if claims.TokenVersion != 1 {
		t.Errorf("expected TokenVersion=1, got %d", claims.TokenVersion)
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	refreshToken, err := auth.CreateRefreshToken(context.Background(), 42, "testuser", 1)
	if err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}

	claims, err := auth.VerifyToken(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}

	if claims.TokenType != "refresh" {
		t.Errorf("expected TokenType=refresh, got %s", claims.TokenType)
	}
}

func TestAuthService_VerifyInvalidToken(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	_, err := auth.VerifyToken(context.Background(), "invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestAuthService_GetUserByUsername(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:             1,
				Username:       "testuser",
				Email:          "test@example.com",
				Nickname:       strPtr("Test"),
				UserType:       "local",
				IsActive:       true,
				TokenVersion:   1,
			},
		},
	}

	auth.userRepo = stub

	user, err := auth.GetUserByUsername(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected ID=1, got %d", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("expected Username=testuser, got %s", user.Username)
	}
}

func TestAuthService_GetUserByEmail(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
	}

	auth.userRepo = stub

	user, err := auth.GetUserByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected Email=test@example.com, got %s", user.Email)
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	auth := NewAuthService(nil, "test-secret-key", time.Hour, 7*24*time.Hour)

	stub := &testutil.StubUserRepositoryWithData{
		Users: map[int]*domain.User{
			1: {
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
	}

	auth.userRepo = stub

	user, err := auth.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected ID=1, got %d", user.ID)
	}
}

// strPtr returns a pointer to the given string, for use in test data.
func strPtr(s string) *string {
	return &s
}
