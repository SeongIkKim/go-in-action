package search

// 기본 검색기를 구현할 defaultMatcher 타입.
// Matcher 인터페이스를 구현한다.
type defaultMatcher struct{}

// init 함수에서는 기본 검색기를 프로그램에 등록한다.
// main 함수 실행 전 실행되는데, defaultMatcher 타입의 값을 생성해서 matchers 맵에 default라는 키로 검색기를 저장해둔다.
func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// Search 함수는 기본 검색기의 동작을 구현한다.
// defaultMatcher 타입에 대한 value receiver를 선언함으로써 해당 메서드는 defaultMatcher 타입의 값이나 포인터에 대해 호출될 수 있다.
// 만약 pointer receiver(m *defaultMatcher)를 선언했다면 해당 메서드는 인터페이스 타입의 값(dm)이 아니라 포인터(&dm)에 대해서만 호출될 수 있다.
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}
