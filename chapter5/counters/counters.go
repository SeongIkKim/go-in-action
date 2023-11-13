package counters

// alertCounter는 패키지 내부에서만 사용되는 unexported identifier이다(소문자로 시작한다)
type alertCounter int

// alertCounter를 생성할 수 있도록 만들어진 팩토리 함수
func New(value int) alertCounter {
	return alertCounter(value)
}
