// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/ellesde/e7api.go/e7"
)

func TestHeroes_GetByID(t *testing.T) {
	client := e7.NewClient()

	tests := []struct {
		in   string
		want string
	}{
		// Normal hero - single name, no exclusive equip, no SC
		{in: "aramintha", want: "Aramintha"},
		// Multi-name hero
		{in: "little-queen-charlotte", want: "Little Queen Charlotte"},
		// Hero with exclusive equip
		{in: "cermia", want: "Cermia"},
		// Hero with specialty change
		{in: "montmorancy", want: "Montmorancy"},
		// Specialty change hero
		{in: "angelic-montmorancy", want: "Angelic Montmorancy"},
	}

	for _, tt := range tests {
		h, _, err := client.Heroes.GetByID(context.Background(), tt.in)
		if err != nil {
			t.Fatalf("Heroes.Get('%s') returned error: %v", tt.in, err)
		}

		if want := tt.in; want != h.UUID {
			t.Errorf("hero.UUID was %q, wanted %q", h.UUID, want)
		}

		if want := tt.want; want != h.Name {
			t.Errorf("hero.Name was %q, wanted %q", h.Name, want)
		}
	}
}
