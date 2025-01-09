package main

import (
	"fmt"
	"log"
)

func main() {
	manifest, err := load("cloakroom.json")
	if err != nil {
		log.Fatalf("Error loading manifest: %v\n", err)
	}

	resolve(manifest)

	if err := restore(manifest); err != nil {
		log.Fatalf("Restore failed: %v\n", err)
	}

	fmt.Println("Restore command completed successfully!")
}
