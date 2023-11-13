// race condition 예시

package main

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	counter int // 모든 고루틴이 값을 증가시키려는 변수
	wg      sync.WaitGroup
)

func main() {
	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	wg.Wait()
	fmt.Println("최종 결과: ", counter) // 4가 나와야하는데, 간혹 2가 나온다.

}

func incCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// counter 변수값 읽기
		value := counter

		// 스레드를 양보하여 큐로 돌아가도록 한다.
		runtime.Gosched() // 작업 도중 스레드를 강제로 양보하여 다른 고루틴이 실행되도록 만든다.

		// 현재 counter 값 증가
		value++

		// 증가된 value 값을 counter 변수에 저장
		counter = value
	}
}

// go build -race <filename> 으로 race condition이 의심되는 코드 부분을 검출할 수 있다.
