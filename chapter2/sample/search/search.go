package search

import (
    "log"
	"sync"
)

// 검색을 처리할 검색기의 매핑 정보를 저장할 맵(map)
// package 수준 변수지만, 소문자로 시작하여 private(외부에서 접근 불가)이다.
// make함수는 런타임에 reference type의 인스턴스인 map을 생성한다.
var matchers = make(map[string]Matcher)

// 검색 로직을 수행할 Run 함수
func Run(searchTerm string) {
	// 검색할 피드의 목록을 조회한다.
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// 언버퍼드 채널을 생성하여 화면에 표시할 검색 결과를 전달 받는다.
	// 제로값 초기화를 할 때에는 var를 사용하지만, 함수 호출이나 다른 초기화 로직을 사용해 변수를 초기화할 때에는 :=를 사용한다.
	results := make(chan *Result)

	// 모든 피드를 처리할 때까지 기다릴 대기 그룹(Wait group)을 설정한다.
	// Go에서는 main 함수가 return되면 프로그램이 종료되므로, 다른 goroutine이 실행되는 동안 main 함수가 종료되지 않도록 하기 위해 사용한다.
	// WaitGroup은 counting sempahore로, Add로 카운트를 늘리고 Done으로 카운트를 줄인다.
	var waitGroup sync.WaitGroup

	// 개별 피드를 처리하는 동안 대기해야 할
	// 고루틴의 갯수를 설정한다.
	waitGroup.Add(len(feeds)) // 세마포어 카운트를 늘린다.

	// 각기 다른 종류의 피드를 처리할 고루틴을 실행한다.
	// for 인덱스(_), 요소복사본(feed) := range 컬렉션(feeds) { }
	for _, feed := range feeds {
		// 검색을 위해 검색기를 조회한다.
		matcher, exists := matchers[feed.Type] // map에서 키 조회 시 두 번째 리턴값으로 존재여부(bool)를 알려준다.
		if !exists {
			matcher = matchers["default"]
		}

		// 검색을 실행하기 위해 고루틴을 실행힌다.
		// go 예약어는 고루틴을 실행시키는데 사용한다.
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done() // 세마포어 카운트를 줄인다.
		}(matcher, feed)
	}

	// 모든 작업이 완료되었는지를 모니터링할 고루틴을 실행한다.
	go func() {
		// 모든 작업이 처리될 때까지 기다린다.
		waitGroup.Wait()

		// Display 함수에게 프로그램을 종료할 수 있음을
		// 알리기 위해 채널을 닫는다.
		close(results)
	}()

	// 검색 결과를 화면에 표시하고
	// 마지막 결과를 표시한 뒤 리턴한다.
	Display(results)
}

// 프로그램에서 사용할 검색기를 등록할 함수를 정의한다.
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "검색기가 이미 등록되었습니다.")
	}

	log.Println("등록 완료:", feedType, " 검색기")
	matchers[feedType] = matcher
}
