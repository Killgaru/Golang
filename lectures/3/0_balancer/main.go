// main
package main

import (
	//"fmt"
	"sync"
	//"time"
)

type RoundRobinBalancer struct {
	servers []int
	mx      sync.Mutex
}

func (r *RoundRobinBalancer) Init(in int) {
	r.servers = make([]int, in)
	return
}

func (r *RoundRobinBalancer) GiveStat() []int {
	return r.servers
}

func (r *RoundRobinBalancer) GiveNode() (out int) {
	r.mx.Lock()
	// var wg sync.WaitGroup
	c := make(chan int)
	// wg.Add(1)
	go r.scanServers(c)
	out = <-c
	// wg.Wait()
	//time.Sleep(time.Millisecond)
	r.mx.Unlock()
	return
}

func (r *RoundRobinBalancer) scanServers(c chan int) {
	// fmt.Println("*****")
	end := len(r.servers) - 1
	for i, v := range r.servers {
		if (i == 0) && (v == r.servers[end]) {
			r.servers[i]++
			c <- i
			break
		}
		if i == end {
			r.servers[end]++
			c <- i
			break
		}
		if i != 0 && v < r.servers[i-1] && v <= r.servers[i+1] {
			r.servers[i]++
			c <- i
			break
		}
	}
	// wg.Done()
	return
}

func main() {

}
