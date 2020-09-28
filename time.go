package null

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/volatiletech/randomize"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
type Time struct {
	Time  time.Time
	Valid bool
	set   bool
}

// NewTime creates a new Time.
func NewTime(t time.Time, valid, set bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
		set:   set,
	}
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, true, true)
}

// TimeFromPtr creates a new Time that will be null if t is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false, true)
	}
	return NewTime(*t, true, true)
}

func (t Time) IsSet() bool {
	return t.set
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) error {
	t.set = true
	if bytes.Equal(data, NullBytes) {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}

	if err := t.Time.UnmarshalJSON(data); err != nil {
		return err
	}

	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Time) UnmarshalText(text []byte) error {
	t.set = true
	if text == nil || len(text) == 0 {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
	t.set = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for an invalid Time's value, for potential future omitempty support.
func (t Time) IsZero() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface{}) error {
	t.set = true
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Randomize for sqlboiler
func (t *Time) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		t.Time = time.Time{}
		t.Valid = false
	} else {
		t.Time = randomize.Date(nextInt)
		t.Valid = true
	}
}
