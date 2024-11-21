package main

import (
	"flag"
	"fmt"
	"gas/pkg/auth"
)

func main() {
	username := flag.String("username", "", "Username for the JWT token")
	phone := flag.String("phone", "", "Phone number for the JWT token")
	notes := flag.String("notes", "", "Notes for the JWT token")
	expiresIn := flag.Int("expiresIn", 0, "Number of days until the token expires (default is 30 days)")
	flag.Parse()

	if *username == "" {
		fmt.Println("Username is required")
		return
	}

	token, err := auth.GenerateJWT(*username, *phone, *notes, *expiresIn)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	fmt.Println("Generated JWT token:")
	fmt.Println(token)
}
