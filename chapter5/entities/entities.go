package entities

type User struct {
	Name  string // Name 필드는 외부로 노출된다(대문자로 시작한다)
	email string // email 필드는 외부로 노출되지 않는다(소문자로 시작한다)
}

type info struct {
	Phone string // Phone 필드는 외부로 노출된다(대문자로 시작한다)
}

type Admin struct {
	User   // exported
	info   // 포함된 타입을 unexported type으로 선언한다.
	Rights int
}
