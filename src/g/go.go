package g

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	wg *sync.WaitGroup
)

func init() {
	wg = new(sync.WaitGroup)
}

func Go(f func()) {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("goroutine recovered: %+v", r)
			}
			wg.Done()
		}()
		f()
	}()
}

func Add(i int) {
	wg.Add(i)
}

func Done() {
	wg.Done()
}

func Wait() {
	wg.Wait()
}
