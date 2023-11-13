package main

import (
	"fmt"
)

type user struct {
	name  string
	email string
}

type notifier interface {
	notify()
}

func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

type admin struct {
	user  // type embedding
	level string
}

func main() {
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	ad.user.notify()
	ad.notify() // inner type's method is promoted, admin에서 직접적으로 user의 notify() 메서드를 호출할 수 있다

	sendNotification(&ad) // user의 notify() 메서드를 호출할 수 있다
}

// notifier 인터페이스를 구현하는 값을 매개변수로 전달받아 알림을 보낸다.
func sendNotification(n notifier) {
	n.notify()
}
