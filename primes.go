package main

import (
	"container/heap"
	"fmt"
)

type sieve struct {
	next int
	ch   chan int
}

type intHeap []sieve

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i].next < h[j].next }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(sieve)) }

func (h *intHeap) Pop() interface{} {
	defer func() { *h = (*h)[0 : h.Len()-1] }()
	return (*h)[h.Len()-1]
}

func makeSieve(ch chan int, prime int) sieve {
	counter := make(chan int, 1)
	count := prime
	go func() {
		count += prime
		counter <- count
	}()
	ch <- prime
	return sieve{<-counter, counter}
}

func primes(ch chan int) {
	h := intHeap{makeSieve(ch, 2), makeSieve(ch, 3)}
	heap.Init(&h)
	c := h[0].next
	for {
		if c != h[0].next {
			heap.Push(&h, makeSieve(ch, c))
		}
		c = h[0].next
		for c == h[0].next {
			_N := heap.Pop(&h)
			n := _N.(sieve)
			c = n.next
			heap.Push(&h, sieve{<-n.ch, n.ch})
		}
		c++
	}
}

func main() {
	ch := make(chan int, 1)
	go primes(ch)
	fmt.Printf("%d\n", <-ch)
	fmt.Printf("%d\n", <-ch)
	fmt.Printf("%d\n", <-ch)
	for <-ch < (1 << 16) {
	}
	fmt.Printf("%d\n", 1<<16)
	fmt.Printf("%d\n", <-ch)
	fmt.Printf("%d\n", 1<<17)
	for {
	}
}
