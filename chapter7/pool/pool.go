// 사용자가 정의한 리소스 집합을 관리하는 패키지
package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

// Pool 구조체는 여러 개의 고루틴에서 안전하게 공유하기 위한 리소스 집합을 관리한다.
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer // 리소스는 io.Closer 인터페이스를 구현해야 한다.
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("풀이 닫혔습니다.")

// New 함수는 풀을 생성하는 팩토리 함수를 받아 size 버퍼 크기의 풀을 생성한다
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("풀의 크기가 너무 작습니다.")
	}

	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		log.Println("공유 리소스 획득:", r)
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("리소스획득:", "새로운 리소스")
		return p.factory()
	}
}

func (p *Pool) Release(r io.Closer) {
	// Mutex를 사용하는 목적은 2가지이다.
	// 1. closed 플래그값을 읽으려는 시점에서 Close 메서드가 실행되어 closed 플래그 값을 쓰는 상황이 일어나지 않도록 하기 위해
	// 2. 닫힌 채널에 리소스를 돌려보내면 패닉이 발생하기 때문에 이를 방지하기 위해
	p.m.Lock()
	defer p.m.Unlock()

	// resources 채널이 닫혔는지 확인하기 위해 closed 플래그를 사용한다.
	if p.closed {
		r.Close()
		return
	}

	select {
	// 해제한 리소스를 다시 리소스 큐에 추가
	case p.resources <- r:
		log.Println("리소스 반환:", "리소스 큐에 반환")

	// 리소스 큐가 가득 찬 경우  리소스를 해제
	default:
		log.Println("리소스 반환:", "리소스 해제")
		r.Close()
	}
}

// 풀을 종료하고 모든 리소스를 해제
// 이 메서드의 전체 코드는 단 하나의 고루틴에 의해서만 실행되어야한다. Release 메서드와 동시에 실행되면 안된다.

func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	// 풀을 닫힌 상태로 전환
	p.closed = true

	// 리소스 해제하기에 앞서 먼저 채널을 닫는다.
	// 그렇지 않으면 deadlock 위험이 있음
	close(p.resources)

	for r := range p.resources {
		r.Close()
	}
}
