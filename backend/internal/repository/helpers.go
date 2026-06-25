package repository

import "time"

// --- Pointer helpers for optional fields ---

func timePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// strPtr is an alias for stringPtr.
func strPtr(s string) *string {
	return stringPtr(s)
}

// intPtr returns nil for zero values.
func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
