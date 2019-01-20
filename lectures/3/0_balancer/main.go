// main
package main

import (
	"fmt"
	//"sync"
	"time"
)

type RoundRobinBalancer struct {
	servers []int
}

func (r *RoundRobinBalancer) Init(in int) {
	r.servers = make([]int, in)
	return
}

func (r *RoundRobinBalancer) GiveStat() []int {
	return r.servers
}

func (r *RoundRobinBalancer) GiveNode() int {
	//var wg sync.WaitGroup
	//wg.Add(1)
	c := make(chan int)
	go r.scanServers(c)
	//wg.Wait()
	time.Sleep(10 * time.Millisecond)
	return <-c
}

func (r *RoundRobinBalancer) scanServers(c chan int) {
	for i, v := range r.servers {
		if i+1 == len(r.servers) {
			if v != r.servers[0] {
				r.servers[i]++
				//wg.Done()
				c <- i + 1
				return
			}
			r.servers[0]++
			//wg.Done()
			c <- i + 1
			return
		}
		if v == r.servers[i+1] {
			r.servers[i]++
			//wg.Done()
			c <- i + 1
			return
		}

	}
	c <- 0
	//wg.Done()
	return
}

func main() {
	fmt.Println("Hello World!")
	balancer := new(RoundRobinBalancer)
	balancer.Init(3)
	fmt.Println(balancer.GiveStat())
	n := balancer.GiveNode()
	fmt.Println(n, balancer.GiveStat())
	n = balancer.GiveNode()
	fmt.Println(n, balancer.GiveStat())

}
