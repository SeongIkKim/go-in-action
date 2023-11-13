package main

import (
	"fmt"

	"../entities"
)

func main() {
	u := entities.User{
		Name: "Bill",
		// email: "bill@email.com", // unknown field 'email' in struct literal of type entities.User
	}

	fmt.Printf("사용자: %v\n", u)
}
