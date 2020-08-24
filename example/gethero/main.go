package main

import (
	"context"
	"fmt"

	"github.com/ellesde/e7api.go/e7"
)

func fetchHero(name string) (*e7.Hero, error) {
	client := e7.NewClient()
	h, _, err := client.Heroes.GetByID(context.Background(), name)
	return h, err
}

func main() {
	var name string
	fmt.Printf("Enter hero name: ")
	fmt.Scanf("%s", &name)

	h, err := fetchHero(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Name: %v\n", h.Name)
	fmt.Printf("Rarity: %v\n", h.Rarity)
	fmt.Printf("Attribute: %v\n", h.Attribute)
	fmt.Printf("Role: %v\n", h.Role)
}
