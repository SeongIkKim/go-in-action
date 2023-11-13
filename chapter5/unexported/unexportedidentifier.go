package main

import (
	"fmt"

	"../counters"
)

func main() {
	// unexported identifier 타입의 변수를 생성하려고 하면 컴파일 에러가 발생한다.
	// counter := counters.alertCounter(10)

	counter := counters.New(10)

	fmt.Printf("카운터: %d\n", counter)
}
