package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main1() {
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", string(hash))
}
