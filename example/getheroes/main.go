package main

import (
	"context"
	"fmt"

	"github.com/ellesde/e7api.go/e7"
)

func fetchHeroes() ([]e7.Hero, error) {
	client := e7.NewClient()
	h, _, err := client.Heroes.List(context.Background())
	return h, err
}

func main() {
	heroes, err := fetchHeroes()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Heroes:\n")
	for _, h := range heroes {
		fmt.Printf("\t%v\n", h.Name)
	}
}
