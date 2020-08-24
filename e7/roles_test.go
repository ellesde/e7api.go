package e7

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		in   Role
		want string
	}{
		{in: Warrior, want: "warrior"},
		{in: Knight, want: "knight"},
		{in: Thief, want: "assassin"},
		{in: Ranger, want: "ranger"},
		{in: Mage, want: "mage"},
		{in: SoulWeaver, want: "manauser"},
	}

	for _, tt := range tests {
		got := tt.in.String()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Role.String mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestRole_MarshalJSON(t *testing.T) {
	tests := []struct {
		in   Role
		want []byte
	}{
		{in: Warrior, want: []byte(`"warrior"`)},
		{in: Knight, want: []byte(`"knight"`)},
		{in: Thief, want: []byte(`"assassin"`)},
		{in: Ranger, want: []byte(`"ranger"`)},
		{in: Mage, want: []byte(`"mage"`)},
		{in: SoulWeaver, want: []byte(`"manauser"`)},
	}

	for _, tt := range tests {
		got, err := tt.in.MarshalJSON()
		if err != nil {
			t.Errorf("Role.MarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Role.MarshalJSON mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestRole_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   []byte
		want Role
	}{
		{in: []byte(`"warrior"`), want: Warrior},
		{in: []byte(`"knight"`), want: Knight},
		{in: []byte(`"assassin"`), want: Thief},
		{in: []byte(`"ranger"`), want: Ranger},
		{in: []byte(`"mage"`), want: Mage},
		{in: []byte(`"manauser"`), want: SoulWeaver},
	}

	for _, tt := range tests {
		r := new(Role)
		err := r.UnmarshalJSON(tt.in)
		if err != nil {
			t.Errorf("Role.UnmarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, *r); diff != "" {
			t.Errorf("Role.UnmarshalJSON mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRole_UnmarshalJSON_unknownRole(t *testing.T) {
	r := new(Role)
	err := r.UnmarshalJSON([]byte(`"test"`))
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if !errors.Is(err, ErrUnknownRole) {
		t.Errorf("expected unknown role error")
	}
}
