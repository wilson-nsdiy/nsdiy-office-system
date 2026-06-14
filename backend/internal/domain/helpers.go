package domain

import "time"

// --- Pointer helpers for optional fields ---

func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// StrPtr is an alias for StringPtr.
func StrPtr(s string) *string {
	return StringPtr(s)
}

// BoolPtr always returns non-nil pointer (booleans have no zero/nil distinction).
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns nil for zero values.
func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// FloatPtr returns nil for zero values.
func FloatPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

// DerefString safely dereferences a string pointer.
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// DerefInt safely dereferences an int pointer.
func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// DerefTime safely dereferences a time pointer.
func DerefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// DerefFloat safely dereferences a float64 pointer.
func DerefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// DerefBool safely dereferences a bool pointer.
func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// IntField returns the int value, or 0 if nil (for Ent SetNillable methods).
func IntField(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
