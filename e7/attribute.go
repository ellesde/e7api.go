package e7

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Attribute represents an Epic Seven hero's attribute.
type Attribute int

// Hero attribute.
const (
	None           = -1
	Fire Attribute = iota
	Ice
	Earth
	Light
	Dark
)

var attributeStrings = map[Attribute]string{
	None:  "none",
	Fire:  "fire",
	Ice:   "ice",
	Earth: "wind",
	Light: "light",
	Dark:  "dark",
}

var attributes = map[string]Attribute{
	"none":  None,
	"fire":  Fire,
	"ice":   Ice,
	"wind":  Earth,
	"light": Light,
	"dark":  Dark,
}

func (a Attribute) String() string {
	return attributeStrings[a]
}

// MarshalJSON marshals a as a quoted JSON string.
func (a Attribute) MarshalJSON() ([]byte, error) {
	buf := writeStringBuffer(attributeStrings[a])
	return buf.Bytes(), nil
}

func writeStringBuffer(s string) *bytes.Buffer {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(s)
	buf.WriteString(`"`)
	return buf
}

// UnmarshalJSON unmarshals a quoted JSON string to a.
func (a *Attribute) UnmarshalJSON(b []byte) error {
	s, err := unmarshalJSON(b)
	if err != nil {
		return err
	}

	val, ok := attributes[s]
	if !ok {
		return ErrUnknownAttribute
	}
	*a = val
	return nil
}

// ErrUnknownAttribute is returned when unmarshalling a quoted JSON string whose
// value is not in the list of defined attributes.
var ErrUnknownAttribute = errors.New("unknown attribute")

func unmarshalJSON(b []byte) (string, error) {
	var s string
	err := json.Unmarshal(b, &s)
	return s, err
}
