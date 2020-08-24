package e7

import "errors"

// Role represents an Epic Seven hero's role.
type Role int

// Hero role.
const (
	Warrior Role = iota
	Knight
	// EpicSevenDB API refers to theifs as assassins.
	Thief
	Ranger
	Mage
	// EpicSevenDB API refers to soul weavers as manausers.
	SoulWeaver
)

var rolesStrings = map[Role]string{
	Warrior:    "warrior",
	Knight:     "knight",
	Thief:      "assassin",
	Ranger:     "ranger",
	Mage:       "mage",
	SoulWeaver: "manauser",
}

var roles = map[string]Role{
	"warrior":  Warrior,
	"knight":   Knight,
	"assassin": Thief,
	"ranger":   Ranger,
	"mage":     Mage,
	"manauser": SoulWeaver,
}

func (r Role) String() string {
	return rolesStrings[r]
}

// MarshalJSON marshals r as a quoted JSON string.
func (r Role) MarshalJSON() ([]byte, error) {
	buf := writeStringBuffer(rolesStrings[r])
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted JSON string to r.
func (r *Role) UnmarshalJSON(b []byte) error {
	s, err := unmarshalJSON(b)
	if err != nil {
		return err
	}

	val, ok := roles[s]
	if !ok {
		return ErrUnknownRole
	}
	*r = val
	return nil
}

// ErrUnknownRole is returned when unmarshalling a quoted JSON string whose
// value is not in the list of defined roles.
var ErrUnknownRole = errors.New("unknown role")
