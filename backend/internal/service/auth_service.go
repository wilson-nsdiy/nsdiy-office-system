package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"oa-nsdiy/backend/internal/domain"
)

// UserRepository defines the interface for user data access required by AuthService.
type UserRepository interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdatePassword(ctx context.Context, id int, salt, hashedPassword string) error
	SetVerificationCode(ctx context.Context, id int, code string, expiresAt time.Time) error
	ClearVerificationCode(ctx context.Context, id int) error
}

type AuthService struct {
	userRepo      UserRepository
	jwtSecret     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(userRepo UserRepository, jwtSecret string, accessExpiry, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     []byte(jwtSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

type UserDTO struct {
	ID             int
	Username       string
	Email          string
	Nickname       *string
	HashedPassword string
	Salt           string
	RoleID         *int
	UserType       string
	IsActive       bool
	TokenVersion   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func userToDTO(u *domain.User) *UserDTO {
	return &UserDTO{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Nickname:       u.Nickname,
		HashedPassword: u.HashedPassword,
		Salt:           u.Salt,
		RoleID:         u.RoleID,
		UserType:       u.UserType,
		IsActive:       u.IsActive,
		TokenVersion:   u.TokenVersion,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*UserDTO, error) {
	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return userToDTO(u), nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return userToDTO(u), nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int) (*UserDTO, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return userToDTO(u), nil
}

type Claims struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	TokenType    string `json:"type"`
	TokenVersion int    `json:"token_version"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateSalt(ctx context.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) HashPassword(ctx context.Context, password, salt string) (string, error) {
	combined := salt + password
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) VerifyPassword(ctx context.Context, password, hashedPassword, salt string) bool {
	combined := salt + password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combined))
	return err == nil
}

func (s *AuthService) CreateAccessToken(ctx context.Context, userID int, username string, tokenVersion int) (string, error) {
	claims := &Claims{
		UserID:       userID,
		Username:     username,
		TokenType:    "access",
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) CreateRefreshToken(ctx context.Context, userID int, username string, tokenVersion int) (string, error) {
	claims := &Claims{
		UserID:       userID,
		Username:     username,
		TokenType:    "refresh",
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GenerateVerificationCode(ctx context.Context) string {
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *AuthService) ValidatePasswordStrength(ctx context.Context, password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

func (s *AuthService) HashToken(ctx context.Context, token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) GenerateApiToken(ctx context.Context) (string, string, string) {
	bytes := make([]byte, 20)
	_, _ = rand.Read(bytes)
	token := hex.EncodeToString(bytes)
	prefix := token[:8]
	hash := s.HashToken(ctx, token)
	return token, prefix, hash
}

// ChangePassword verifies the old password, validates the new one, and updates.
// TokenVersion is automatically incremented by UpdatePassword.
func (s *AuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !s.VerifyPassword(ctx, oldPassword, u.HashedPassword, u.Salt) {
		return fmt.Errorf("old password is incorrect")
	}

	if !s.ValidatePasswordStrength(ctx, newPassword) {
		return fmt.Errorf("new password does not meet strength requirements")
	}

	salt, err := s.GenerateSalt(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	hashedPassword, err := s.HashPassword(ctx, newPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, salt, hashedPassword)
}

// RequestResetPassword generates a verification code for password reset and stores it on the user.
// Returns the code (for testing; in production this would be sent via email).
func (s *AuthService) RequestResetPassword(ctx context.Context, email string) (code string, expiresAt time.Time, err error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("user with email %s not found", email)
	}

	code = s.GenerateVerificationCode(ctx)
	expiresAt = time.Now().Add(30 * time.Minute)

	if err := s.userRepo.SetVerificationCode(ctx, u.ID, code, expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to store verification code: %w", err)
	}

	return code, expiresAt, nil
}

// ConfirmResetPassword verifies the reset code and changes the password.
func (s *AuthService) ConfirmResetPassword(ctx context.Context, email, code, newPassword string) error {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if u.VerificationCode == nil || *u.VerificationCode != code {
		return fmt.Errorf("invalid verification code")
	}

	if u.VerificationCodeExpiresAt == nil || time.Now().After(*u.VerificationCodeExpiresAt) {
		return fmt.Errorf("verification code has expired")
	}

	if !s.ValidatePasswordStrength(ctx, newPassword) {
		return fmt.Errorf("password does not meet strength requirements")
	}

	salt, err := s.GenerateSalt(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	hashedPassword, err := s.HashPassword(ctx, newPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, u.ID, salt, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Clear the verification code after successful password reset
	_ = s.userRepo.ClearVerificationCode(ctx, u.ID)

	return nil
}
