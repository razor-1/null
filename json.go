package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/volatiletech/null/v9/convert"
	"github.com/volatiletech/randomize"
)

// JSON is a nullable []byte.
type JSON struct {
	JSON  []byte
	Valid bool
	Set   bool
}

// NewJSON creates a new JSON
func NewJSON(b []byte, valid, set bool) JSON {
	return JSON{
		JSON:  b,
		Valid: valid,
		Set:   set,
	}
}

// JSONFrom creates a new JSON that will be invalid if nil.
func JSONFrom(b []byte) JSON {
	return NewJSON(b, b != nil, true)
}

// JSONFromPtr creates a new JSON that will be invalid if nil.
func JSONFromPtr(b *[]byte) JSON {
	if b == nil {
		return NewJSON(nil, false, true)
	}
	n := NewJSON(*b, true, true)
	return n
}

func (j JSON) IsSet() bool {
	return j.Set
}

// Unmarshal will unmarshal your JSON stored in
// your JSON object and store the result in the
// value pointed to by dest.
func (j JSON) Unmarshal(dest interface{}) error {
	if dest == nil {
		return errors.New("destination is nil, not a valid pointer to an object")
	}

	// Call our implementation of
	// JSON MarshalJSON through json.Marshal
	// to get the value of the JSON object
	res, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return json.Unmarshal(res, dest)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JSON) UnmarshalJSON(data []byte) error {
	j.Set = true
	if data == nil {
		return fmt.Errorf("json: cannot unmarshal nil into Go value of type null.JSON")
	}

	if bytes.Equal(data, NullBytes) {
		j.JSON = NullBytes
		j.Valid = false
		return nil
	}

	j.Valid = true
	j.JSON = make([]byte, len(data))
	copy(j.JSON, data)

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (j *JSON) UnmarshalText(text []byte) error {
	j.Set = true
	if text == nil || len(text) == 0 {
		j.JSON = nil
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], text...)
		j.Valid = true
	}

	return nil
}

// Marshal will marshal the passed in object,
// and store it in the JSON member on the JSON object.
func (j *JSON) Marshal(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Call our implementation of
	// JSON UnmarshalJSON through json.Unmarshal
	// to Set the result to the JSON object
	return json.Unmarshal(res, j)
}

// MarshalJSON implements json.Marshaler.
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j.JSON) == 0 || j.JSON == nil {
		return NullBytes, nil
	}
	return j.JSON, nil
}

// MarshalText implements encoding.TextMarshaler.
func (j JSON) MarshalText() ([]byte, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// SetValid changes this JSON's value and also sets it to be non-null.
func (j *JSON) SetValid(n []byte) {
	j.JSON = n
	j.Valid = true
	j.Set = true
}

// Ptr returns a pointer to this JSON's value, or a nil pointer if this JSON is null.
func (j JSON) Ptr() *[]byte {
	if !j.Valid {
		return nil
	}
	return &j.JSON
}

// IsZero returns true for null or zero JSON's, for future omitempty support (Go 1.4?)
func (j JSON) IsZero() bool {
	return !j.Valid
}

// Scan implements the Scanner interface.
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		j.JSON, j.Valid, j.Set = nil, false, false
		return nil
	}
	j.Valid, j.Set = true, true
	return convert.ConvertAssign(&j.JSON, value)
}

// Value implements the driver Valuer interface.
func (j JSON) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// Randomize for sqlboiler
func (j *JSON) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	j.JSON = []byte(`"` + randomize.Str(nextInt, 1) + `"`)
	j.Valid = true
}
