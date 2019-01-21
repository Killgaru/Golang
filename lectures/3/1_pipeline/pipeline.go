package pipeline

import (
	"sync"
)

type job func(in, out chan interface{})

func (j job) oneJob(in chan interface{}, wg *sync.WaitGroup) chan interface{} {
	out := make(chan interface{})
	out2 := make(chan interface{})
	quit := make(chan bool)
	go func() {
		j(in, out)
		quit <- true
	}()
	go func() {
		var val []interface{}
		for {
			select {
			case v := <-out:
				val = append(val, v)
			case <-quit:
				wg.Done()
				for _, v := range val {
					out2 <- v
				}
				close(out2)
				return
			}
		}
	}()
	return out2
}

func Pipe(funcs ...job) {
	var wg sync.WaitGroup
	var inni chan interface{}
	for _, v := range funcs {
		wg.Add(1)
		inni = v.oneJob(inni, &wg)
		wg.Wait()
	}
	return
}
