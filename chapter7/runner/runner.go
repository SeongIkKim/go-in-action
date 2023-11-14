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
	complete chan error // 처리가 종료되었음을 알리기 위한 채널

	timeout <-chan time.Time // 처리 시간을 초과했음을 알리기 위한 채널

	tasks []func(int) // 인덱스 순서로 처리될 작업 목록을 저장하기 위한 슬라이스
}

// timeout 채널에서 값을 수신하면 ErrTimeout을 리턴
var ErrTimeout = errors.New("시간 초과")

// interrupt 채널에서 값을 수신하면 ErrInterrupt를 리턴
var ErrInterrupt = errors.New("운영체제 인터럽트")

// 실행할 Runner 타입 값을 리턴하는 함수
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d), // 지정된 시간이 지나면 값을 수신하는 채널 생성
	}
}

// Runner 타입에 작업을 추가하는 메서드
// 작업은 int형 Id를 매변수로 전달받는 함수다.
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
	// 작업 완료 신호를 수신한 경우
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

	// 작업을 계속 실행
	default:
		return false
	}
}
