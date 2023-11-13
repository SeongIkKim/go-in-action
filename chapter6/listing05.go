package main

import (
	"fmt"
	"runtime"
)

// 논리 프로세서 1개, 고루틴 2개
func main() {
	// 스케줄러가 사용할 수 있는 논리 프로세서의 최대 개수를 설정한다.
	cpuNum := runtime.NumCPU()
	fmt.Println("cpuNum:", cpuNum) // cpuNum: 12 (내 컴퓨터의 논리 프로세서 개수)
}
