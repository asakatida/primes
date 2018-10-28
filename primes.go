package main

import (
	"flag"
	"log"
)

func runCounter(counter chan<- uint, prime uint) {
	for count := prime * 2; ; count += prime {
		counter <- count
	}
}

func newCounter(prime uint) <-chan uint {
	counter := make(chan uint)
	go runCounter(counter, prime)
	return counter
}

func stepCounter(h map[uint]<-chan uint, counter <-chan uint) {
	c := <-counter
	if h[c], counter = counter, h[c]; counter != nil {
		go stepCounter(h, counter)
	}
}

func insertPrime(primes chan<- uint, h map[uint]<-chan uint, gC uint) {
	primes <- gC
	stepCounter(h, newCounter(gC))
}

func primeStep(primes chan<- uint, h map[uint]<-chan uint, gC uint) {
	counter := h[gC]
	if counter == nil {
		insertPrime(primes, h, gC)
	} else {
		stepCounter(h, counter)
		delete(h, gC)
	}
}

func primeLoop(primes chan<- uint) {
	primes <- 2
	h := make(map[uint]<-chan uint)
	for gC := 3; ; gC += 2 {
		primeStep(primes, h, uint(gC))
	}
}

func makePrimes() <-chan uint {
	primes := make(chan uint)
	go primeLoop(primes)
	return primes
}

func main() {
	var start uint
	var end uint
	flag.UintVar(&start, "start", 0, "usage")
	flag.UintVar(&end, "end", 10, "usage")
	flag.Parse()
	primes := makePrimes()
	for {
		prime := <-primes
		if prime > end {
			return
		}
		if prime >= start {
			log.Println(prime)
		}
	}
}
