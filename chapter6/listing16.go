// access sychronization이 필요한 코드에 mutex를 이용하여 critical section을 생성해서 race condition 방지하는 예제
package main

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	counter int64 // 모든 고루틴이 값을 증가시키려는 변수
	wg      sync.WaitGroup
	mutex   sync.Mutex // critical section 설정시 사용할 mutex
)

func main() {
	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	wg.Wait()
	fmt.Println("최종 결과: %d\n", counter) // always 4
}

// 패키지 수준에 정의된 counter 변수 값을 mutex를 이용하여 안전하게 증가시킨다.
func incCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 이 critical section에는 하나의 goroutine만이 접근할 수 있다.
		mutex.Lock() // critical section 시작
		{            // 괄호는 가독성을 위한 것이지 반드시 필요한 것은 아니다. mutex.Unlock 전까지 critical section을 형성한다.
			// counter 변수값 읽기
			value := counter

			// 스레드를 양보하여 큐로 돌아가도록 한다.
			runtime.Gosched() // 작업 도중 스레드를 강제로 양보하여 다른 고루틴이 실행되도록 만든다.

			// 현재 counter 값 증가
			value++

			// 증가된 value 값을 counter 변수에 저장
			counter = value
		}
		mutex.Unlock() // critical section 종료
		// 이제 다른 고루틴이 critical section에 접근할 수 있다.
	}
}
