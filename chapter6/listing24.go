// 버퍼가 있는 채널을 이용해 미리 정해진 고루틴의 개수만큼 다중 작업을 수행하는 예제
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numberGoroutines = 4  // 실행할 고루틴의 개수
	taskLoad         = 10 // 처리할 작업의 개수
)

var wg sync.WaitGroup

// Go runtime이 다른 코드를 실행하기에 앞서 패키지 초기화를 수행하는 함수
func init() {
	rand.Seed(time.Now().UnixNano()) // 난수 생성기 초기화
}

func main() {
	tasks := make(chan string, taskLoad) // 버퍼 크기가 taskLoad인 채널을 생성한다.

	// 작업을 처리할 고루틴을 실행한다.
	wg.Add(numberGoroutines)
	for gr := 1; gr <= numberGoroutines; gr++ {
		go worker(tasks, gr)
	}

	// 실행할 작업을 추가한다.
	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("작업[%d]", post)
	}

	// 작업을 모두 추가하고 나면 채널을 닫는다.(producer가 더이상 작업을 추가하지 않을 것을 알리는 것이라고 생각하면 된다.)
	// 채널이 닫히더라도 각 고루틴은 더이상 받을 값이 없을 때 까지 채널에서 계속 값을 받는다. 이렇게 함으로써 채널이 닫히더라도 채널 내의 값이 유실되지 않는다.
	close(tasks)

	wg.Wait()
}

// 버퍼가 있는 채널에서 수행할 작업을 가져가는 고루틴
func worker(tasks chan string, worker int) {
	defer wg.Done()

	for {
		// 작업이 할당될 때까지 대기한다.
		// 채널이 닫혀있더라도 채널 내에 남은 값이 있을때까지 계속 값을 받는다.
		task, ok := <-tasks // 채널 내에 값이 남아있다면 닫힘 여부와 관계엾이 ok는 계속 true다. 채널이 비어있거나 비어있으면서 닫혀있는 경우 ok는 false가 된다.
		if !ok {
			// 채널이 닫힌 경우
			fmt.Printf("작업자<%d> : 종료합니다.\n", worker)
			return
		}

		// 작업을 시작하는 메시지를 출력한다.
		fmt.Printf("작업자<%d> : 작업 시작: %s\n", worker, task)

		// 작업을 처리하는 것을 흉내내기 위해 임의의 시간을 대기한다.
		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		// 작업을 종료하는 메시지를 출력한다.
		fmt.Printf("작업자<%d> : 작업 완료: %s\n", worker, task)
	}

}

/*
작업자<2> : 작업 시작: 작업[1]
작업자<4> : 작업 시작: 작업[2]
작업자<1> : 작업 시작: 작업[3]
작업자<3> : 작업 시작: 작업[4]
작업자<4> : 작업 완료: 작업[2]
작업자<4> : 작업 시작: 작업[5]
작업자<3> : 작업 완료: 작업[4]
작업자<3> : 작업 시작: 작업[6]
작업자<2> : 작업 완료: 작업[1]
작업자<2> : 작업 시작: 작업[7]
작업자<1> : 작업 완료: 작업[3]
작업자<1> : 작업 시작: 작업[8]
작업자<1> : 작업 완료: 작업[8]
작업자<1> : 작업 시작: 작업[9]
작업자<4> : 작업 완료: 작업[5]
작업자<4> : 작업 시작: 작업[10]
작업자<2> : 작업 완료: 작업[7]
작업자<2> : 종료합니다.
작업자<4> : 작업 완료: 작업[10]
작업자<4> : 종료합니다.
작업자<3> : 작업 완료: 작업[6]
작업자<3> : 종료합니다.
작업자<1> : 작업 완료: 작업[9]
작업자<1> : 종료합니다.
*/
