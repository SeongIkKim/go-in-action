// pool 패키지를 활용하여 데이터베이스 연결 풀을 생성하고 활용하는 예제
package main

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"../../pool"
)

const (
	maxGoroutines   = 25 // 풀에 저장할 리소스의 개수
	pooledResources = 2  // 풀이 관리할 리소스의 개수
)

type dbConnection struct {
	ID int32
}

// dbConnection 타입이 풀에 의해 관리되도록 하기 위해 io.Closer 인터페이스를 구현한다.
func (dbConn *dbConnection) Close() error {
	log.Println("닫힘: 데이터베이스 연결", dbConn.ID)
	return nil
}

var idCounter int32

func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("생성: 새 데이터베이스 연결", id)

	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	p, err := pool.New(createConnection, pooledResources)
	if err != nil {
		log.Println(err)
	}

	for query := 0; query < maxGoroutines; query++ {
		// 각 고루틴은 풀에서 리소스를 획득한 후, 쿼리를 실행한다.
		// 쿼리가 실행되면 리소스를 풀로 돌려보낸다.
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()

	log.Println("프로그램을 종료합니다.")
	p.Close()
}

func performQueries(query int, p *pool.Pool) {
	// 풀에서 dbConnection 리소스 획득
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}

	defer p.Release(conn)

	// 쿼리문이 실행되는 것처럼 흉내낸다.
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("질의: QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)
}
