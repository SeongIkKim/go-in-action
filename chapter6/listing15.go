package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	shutdown int64 // 실행중인 고루틴들의 종료 신호로 사용될 플래그
	wg       sync.WaitGroup
)

func main() {
	wg.Add(2)

	go doWork("A")
	go doWork("B")

	time.Sleep(1 * time.Second) // 고루틴이 실행될 시간을 할애한다.

	// 종료 신호 플래그를 설정한다.
	fmt.Println("프로그램 종료")
	atomic.StoreInt64(&shutdown, 1)

	wg.Wait()

}

func doWork(name string) {
	defer wg.Done()

	for {
		fmt.Printf("작업자: %s 작업 진행중 ...\n", name)
		time.Sleep(250 * time.Millisecond)

		if atomic.LoadInt64(&shutdown) == 1 { // shutdown 플래그가 1이면 작업을 종료한다.
			fmt.Printf("작업자: %s 종료\n", name)
			break
		}
	}
}

/*
작업자: B 작업 진행중 ...
작업자: A 작업 진행중 ...
작업자: A 작업 진행중 ...
작업자: B 작업 진행중 ...
작업자: B 작업 진행중 ...
작업자: A 작업 진행중 ...
작업자: A 작업 진행중 ...
작업자: B 작업 진행중 ...
프로그램 종료
작업자: B 종료
작업자: A 종료
*/
// main 함수에서 storeInt64 호출 이후 각 goroutine에서 LoadInt64 호출을 통해 shutdown 플래그를 확인하고 종료한다.
