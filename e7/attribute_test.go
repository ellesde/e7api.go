package e7

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAttribute_String(t *testing.T) {
	tests := []struct {
		in   Attribute
		want string
	}{
		{in: None, want: "none"},
		{in: Fire, want: "fire"},
		{in: Ice, want: "ice"},
		{in: Earth, want: "wind"},
		{in: Light, want: "light"},
		{in: Dark, want: "dark"},
	}

	for _, tt := range tests {
		got := tt.in.String()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Attribute.String mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestAttribute_MarshalJSON(t *testing.T) {
	tests := []struct {
		in   Attribute
		want []byte
	}{
		{in: None, want: []byte(`"none"`)},
		{in: Fire, want: []byte(`"fire"`)},
		{in: Ice, want: []byte(`"ice"`)},
		{in: Earth, want: []byte(`"wind"`)},
		{in: Light, want: []byte(`"light"`)},
		{in: Dark, want: []byte(`"dark"`)},
	}

	for _, tt := range tests {
		got, err := tt.in.MarshalJSON()
		if err != nil {
			t.Errorf("Attribute.MarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Attribute.String mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestAttribute_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   []byte
		want Attribute
	}{
		{in: []byte(`"none"`), want: None},
		{in: []byte(`"fire"`), want: Fire},
		{in: []byte(`"ice"`), want: Ice},
		{in: []byte(`"wind"`), want: Earth},
		{in: []byte(`"light"`), want: Light},
		{in: []byte(`"dark"`), want: Dark},
	}

	for _, tt := range tests {
		a := new(Attribute)
		err := a.UnmarshalJSON(tt.in)
		if err != nil {
			t.Errorf("Attribute.UnmarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, *a); diff != "" {
			t.Errorf("Attribute.UnmarshalJSON mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestAttribute_UnmarshalJSON_unknownAttribute(t *testing.T) {
	a := new(Attribute)
	err := a.UnmarshalJSON([]byte(`"test"`))
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if !errors.Is(err, ErrUnknownAttribute) {
		t.Errorf("expected unknown attribute error")
	}
}
