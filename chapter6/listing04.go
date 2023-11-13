package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

// 논리 프로세서 1개, 고루틴 2개
func main() {
	// 스케줄러가 사용할 수 있는 논리 프로세서의 최대 개수를 설정한다.
	runtime.GOMAXPROCS(1) // 논리프로세서를 하나만 사용하기

	wg.Add(2)

	fmt.Println("고루틴을 실행합니다.")

	go printPrime("A")
	go printPrime("B")

	// 고루틴이 종료될 때까지 대기한다
	fmt.Println("대기 중...")
	wg.Wait()

	fmt.Println("\n프로그램을 종료합니다.")
}

// 1부터 5000까지 소수를 출력한다.
func printPrime(prefix string) {
	defer wg.Done()

next:
	for outer := 2; outer < 5000; outer++ {
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s: %d\n", prefix, outer)
	}
	fmt.Println("종료: ", prefix)
}

/*
고루틴을 실행합니다.
대기 중...
A: 2
A: 3
...
B: 2
B: 3
...
A: 4999
..

*/
// 스케줄러에 의해 고루틴이 교체되면서 실행되었다.
