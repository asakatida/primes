package main

import "log"

func runCounter(counter chan<- int, prime int) {
	for count := prime * 2; ; count += prime {
		counter <- count
	}
}

func newCounter(prime int) <-chan int {
	counter := make(chan int)
	go runCounter(counter, prime)
	return counter
}

func stepCounter(h map[int]<-chan int, counter <-chan int) {
	c := <-counter
	if h[c], counter = counter, h[c]; counter != nil {
		go stepCounter(h, counter)
	}
}

func insertPrime(primes chan<- int, h map[int]<-chan int, gC int) {
	primes <- gC
	stepCounter(h, newCounter(gC))
}

func primeStep(primes chan<- int, h map[int]<-chan int, gC int) {
	counter := h[gC]
	if counter == nil {
		insertPrime(primes, h, gC)
	} else {
		stepCounter(h, counter)
		delete(h, gC)
	}
}

func primeLoop(primes chan<- int) {
	primes <- 2
	h := make(map[int]<-chan int)
	for gC := 3; ; gC += 2 {
		primeStep(primes, h, gC)
	}
}

func makePrimes() <-chan int {
	primes := make(chan int)
	go primeLoop(primes)
	return primes
}

func main() {
	primes := makePrimes()
	for {
		log.Println(<-primes)
	}
}
