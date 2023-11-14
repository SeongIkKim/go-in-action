// 프로세스의 실행 및 수명주기를 관리하는 패키지
// 크론 잡을 수행하거나 백그라운드 작업 프로세스를 예약 실행할 때 사용하는 패턴
package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// Runner는 주어진 타임아웃 시간동안 일련의 작업을 수행한다.
// 운영체제 인터럽트에 의해 실행이 종료된다.
type Runner struct {
	interrupt chan os.Signal // 운영체제로부터 전달되는 인터럽트 신호를 수신하는 채널

	// 특이하게도 error interface로 complete를 처리한다
	// 에러가 발생했다면 error 인터페이스 값을 보내고, 에러가 없다면 (정상적으로면 완료했다면) nil 값을 보낸다.
	complete chan error // 처리가 종료되었음을 알리기 위한 채널

	// 모든 작업을 처리하기 위해 사용한 시간을 관리하는 채널
	// time.Time 값을 받지 못하면 프로그램은 자신이 하던 작업을 마무리하고 종료된다.
	timeout <-chan time.Time

	// 이 슬라이스는 int를 매개변수로 받는 (고루틴으로 실행 될) 함수들을 저장한다.
	// int 매개변수는 작업이 수행될 순서를 의미한다.
	tasks []func(int)
}

// timeout 채널에서 값을 수신하면 ErrTimeout을 리턴
var ErrTimeout = errors.New("시간 초과")

// interrupt 채널에서 값을 수신하면 ErrInterrupt를 리턴
var ErrInterrupt = errors.New("운영체제 인터럽트")

// 실행할 Runner 타입 값을 리턴하는 팩토리 함수
func New(d time.Duration) *Runner {
	return &Runner{
		// 각 필드 채널 초기화
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d), // 지정된 시간이 지나면 값을 수신하는 채널 생성
		// tasks는 제로 값이 nil 슬라이스이므로 따로 초기화해주지 않아도 된다.
	}
}

// Runner 타입에 가변길이의 func 매개변수들을 전달받아 runner의 tasks 필드 슬라이스에 추가하는 메서드
// 이 때 각 func는 아무것도 리턴하지 않는 함수여야 한다.
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// 저장된 모든 작업을 실행하고 채널 이벤트를 관찰한다.
func (r *Runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt) // 모든 종류의 인터럽트를 수신한다.

	// 각각의 작업을 각기 다른 고루틴에서 실행
	go func() {
		r.complete <- r.run()
	}()

	select {
	// 작업 완료 신호를 수신한 경우 혹은 도중에 에러가 발생했다는 신호를 수신한 경우
	case err := <-r.complete:
		return err

	// 작업 시간 초과 신호를 수신한 경우
	case <-r.timeout:
		return ErrTimeout
	}
}

// 개별 작업을 실행하는 메서드
func (r *Runner) run() error {
	for id, task := range r.tasks {
		// OS로부터 인터럽트 신호를 수신했는지 확인한다.
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		// 작업을 실행
		task(id)
	}

	return nil
}

// 인터럽트 신호가 수신되었는지 확인하는 메서드
func (r *Runner) gotInterrupt() bool {
	select {
	// 인터럽트 이벤트 발생한 경우
	case <-r.interrupt:
		// 이후에 발생하는 인터럽트 신호를 더이상 수신하지 않도록 한다
		signal.Stop(r.interrupt)
		return true

	// 인터럽트가 발생하지 않았다면 작업을 계속 실행한다.
	// interrupt에서 값이 들어올때까지 기다리는 것이 아니라 바로 다음으로 넘어가기 때문에 non-blocking이다.
	default:
		return false
	}
}
