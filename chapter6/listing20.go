// 두개의 고루틴을 이용하여 테니스 경기를 모방하는 예제
// 버퍼가 없는 채널 사용
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano()) // 난수 생성기 초기화
}

func main() {
	court := make(chan int) // 버퍼가 없는 채널을 생성한다.(버퍼 크기 명시 X)

	wg.Add(2)

	// 두 명의 플레이어를 등록한다.
	// 이 시점부터 두 goroutine은 공을 기다리기 위해 잠금상태가 된다.
	go player("Nadal", court)
	go player("Djokovic", court)

	// 경기를 시작한다.
	court <- 1 // 첫번째 공을 보낸다

	// 경기가 끝날때까지 기다린다.
	wg.Wait()
}

// 테니스 플레이어를 모방하는 함수
func player(name string, court chan int) {
	defer wg.Done()

	for {
		// 공이 되돌아올 때까지 기다린다.
		ball, ok := <-court // 채널에서 값을 읽어오기 위해 대기한다. 이 작업으로 고루틴은 잠금상태에 놓인다.
		if !ok {
			// 채널이 닫혔으면 이긴 것이다.
			fmt.Printf("%s 선수가 승리했습니다.\n", name)
			return
		}

		// 랜덤 값을 이용해 공을 받아치지 못했는지 확인한다.
		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("%s 선수가 공을 받아치지 못했습니다.\n", name)

			// 채널을 닫아 현재 플레이어가 패배했음을 알린다.
			close(court)
			return
		}

		// 선수가 공을 받아친 횟수를 출력하고 공을 상대방에게 보낸다.
		fmt.Printf("%s 선수가 %d 번째 공을 받아쳤습니다.\n", name, ball)
		ball++

		// 공을 상대방에게 보낸다.
		court <- ball
	}

}

/*
Djokovic 선수가 1 번째 공을 받아쳤습니다.
Nadal 선수가 2 번째 공을 받아쳤습니다.
Djokovic 선수가 3 번째 공을 받아쳤습니다.
Nadal 선수가 4 번째 공을 받아쳤습니다.
Djokovic 선수가 5 번째 공을 받아쳤습니다.
Nadal 선수가 6 번째 공을 받아쳤습니다.
Djokovic 선수가 7 번째 공을 받아쳤습니다.
Nadal 선수가 공을 받아치지 못했습니다.
Djokovic 선수가 승리했습니다.
*/
