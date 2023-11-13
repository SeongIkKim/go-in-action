// race condition 예시

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	counter int64 // 모든 고루틴이 값을 증가시키려는 변수
	wg      sync.WaitGroup
)

func main() {
	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	wg.Wait()
	fmt.Println("최종 결과: ", counter) // always 4

}

func incCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		atomic.AddInt64(&counter, 1) // counter 변수에 대한 원자적 연산을 수행하여 race condition을 방지한다.

		runtime.Gosched() // 현재 고루틴이 스레드를 양보하여 다른 고루틴이 실행되도록 한다.
	}
}

// go build -race <filename> 으로 race condition이 의심되는 코드 부분을 검출할 수 있다.
