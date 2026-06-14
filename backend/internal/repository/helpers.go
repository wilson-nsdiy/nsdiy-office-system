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

// boolPtr always returns non-nil pointer (booleans have no zero/nil distinction).
func boolPtr(b bool) *bool {
	return &b
}

// intPtr returns nil for zero values.
func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// floatPtr returns nil for zero values.
func floatPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

// derefString safely dereferences a string pointer.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefInt safely dereferences an int pointer.
func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// derefTime safely dereferences a time pointer.
func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// derefFloat safely dereferences a float64 pointer.
func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// derefBool safely dereferences a bool pointer.
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// intField returns the int value, or 0 if nil (for Ent SetNillable methods).
func intField(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
