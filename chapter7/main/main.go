// 프로그램이 지정된 시간보다 오래 실행중이면 자동으로 종료하기 위해 채널을 활용하는 방법을 소개하는 예제
package main

import (
	"log"
	"os"
	"time"

	"../runner"
)

const timeout = 3 * time.Second

func main() {
	log.Println("작업을 시작합니다.")

	// 실행 시간을 제한하여 프로그램을 실행한다.
	r := runner.New(timeout)

	// 수행할 작업 등록
	r.Add(createTask(), createTask(), createTask())

	// 작업 실행 후 결과 처리
	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Println("지정된 시간을 초과했습니다.")
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println("운영체제 인터럽트가 발생했습니다.")
			os.Exit(2)
		}
	}

	log.Println("프로그램을 종료합니다.")
	// 정상적으로 종료하는 것은 os.Exit(0)를 수행하는것이나 다름없다.
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("프로세서-작업 #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
