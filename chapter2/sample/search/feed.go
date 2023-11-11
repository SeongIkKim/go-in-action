package search

import (
	"encoding/json"
	"os"
)

// 소문자로 시작하는 private 상수, 따라서 search 외부의 패키지에서는 접근할 수 없다.
const dataFile = "data/data.json"

// 피드를 처리할 정보를 표현하는 구조체 (외부로 노출, public)
type Feed struct {
	Name string `json:"site"` // json 태그를 이용해 json 파일의 필드와 매핑한다.
	URI  string `json:"link"`
	Type string `json:"type"`
}

// RetrieveFeeds 함수는 피드 데이터 파일을 읽어 구조체로 변환한다.
// 파라미터 없고 리턴값이 2개, 
func RetrieveFeeds() ([]*Feed, error) {
	// 파일을 연다.
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	// defer 함수를 이용해 이 함수가 리턴될 때
	// 앞서 열어둔 파일이 닫히도록 한다.
	// defer 예약어는 함수가 리턴된 직후에 실행될 작업을 예약한다.(panic 상태에 빠지더라도 반드시 실행된다.)
	defer file.Close()

	// 파일을 읽어 Feed 구조체의 포인터의
	// 슬라이스로 변환한다.
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	// 호출 함수가 오류를 처리할 수 있으므로 오류 처리는 하지 않는다.
	return feeds, err
}
