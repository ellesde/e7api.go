package e7

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHeroesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/hero/h", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"results": [
				{
					"_id": "h",
					"id": "1",
					"name": "H"
				}
			],
			"meta": {
				"requestDate": "date",
				"apiVersion": "1"
			}
		}`)
	})

	ctx := context.Background()
	got, _, err := client.Heroes.GetByID(ctx, "h")
	if err != nil {
		t.Errorf("Heroes.Get returned error: %v", err)
	}

	want := &Hero{
		UUID: "h",
		ID:   "1",
		Name: "H",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Heroes.Get mismatch (-want +got):\n%s", diff)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Heroes.GetByID(ctx, "h")
	if got != nil {
		t.Errorf("client.BaseURL.Path='' GetByID = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' GetByID resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' GetByID err = nil, want error")
	}
}
