package main

import (
	"fmt"

	"../entities"
)

func main() {
	u := entities.Admin{
		Rights: 10,
	}

	u.Name = "Bill"
	// u.email = "bill@email.com" // unknown field 'email' in struct literal of type entities.User
	u.Phone = "awesome-phone-num" // unexported inner type의 exported field에는 접근할 수 있다. inner type identifier가 promoted되었기 때문이다.

	fmt.Printf("사용자: %v\n", u)
}
