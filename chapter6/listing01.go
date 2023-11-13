package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 논리 프로세서 1개, 고루틴 2개
func main() {
	// 스케줄러가 사용할 수 있는 논리 프로세서의 최대 개수를 설정한다.
	runtime.GOMAXPROCS(1) // 논리프로세서를 하나만 사용하기

	// wg는 프로그램의 종료를 대기하기 위해 사용한다.
	// 각각의 고루틴마다 하나씩, 총 두 개의 카운터를 추가한다.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("고루틴을 실행합니다.")

	// 익명함수를 선언하고 고루틴을 생성한다.
	go func() {
		// main 함수에게 종료를 알리기 위한 Done 함수 호출을 예약한다.
		defer wg.Done()

		// 고루틴 1의 작업
		for count := 0; count < 3; count++ {
			for char := 'a'; char < 'a'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
		fmt.Println("\n")
	}()

	// 익명함수를 선언하고 고루틴을 생성한다
	go func() {
		defer wg.Done()

		// 고루틴 2의 작업
		for count := 0; count < 3; count++ {
			for char := 'A'; char < 'A'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
		fmt.Println("\n")
	}()

	// 고루틴이 종료될 때까지 대기한다
	fmt.Println("대기 중...")
	wg.Wait()

	fmt.Println("\n프로그램을 종료합니다.")
}

/*
고루틴을 실행합니다.
대기 중...
A B C D E F G H I J K L M N O P Q R S T U V W X Y Z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z

a b c d e f g h i j k l m n o p q r s t u v w x y z a b c d e f g h i j k l m n o p q r s t u v w x y z a b c d e f g h i j k l m n o p q r s t u v w x y z


프로그램을 종료합니다.

*/
// 짧게 끝나는 작업이라 도중에 스케줄링이 되지않아 고루틴이 순차적으로 실행된 것처럼 보인다.
