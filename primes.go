package main

import (
	"flag"
	"fmt"
)

func generate(ch chan<- uint, start, end uint) {
	defer close(ch)
	if end == 0 {
		for i := uint(3); i < start<<1; i += 2 {
			ch <- i
		}
	} else {
		for i := uint(3); i < end; i += 2 {
			ch <- i
		}
	}
}

func sieve(src <-chan uint, dst chan<- uint, prime uint) {
	defer close(dst)
	for i := range src {
		if i%prime != 0 {
			dst <- i
		}
	}
}

func sieveChannel(src <-chan uint, start, prime uint) <-chan uint {
	src_prime := make(chan uint, 1)
	go sieve(src, src_prime, prime)
	return src_prime
}

func nextPrime(src <-chan uint, prime uint) uint {
	for i := range src {
		return i
	}
	return prime
}

func nextPrimeSend(src <-chan uint, dst chan<- uint, prime uint) uint {
	for i := range src {
		dst <- i
		return i
	}
	return prime
}

func filterRoutine5(src <-chan uint, dst chan<- uint, start, end, prime uint) {
	defer close(dst)
	for last_prime := uint(0); prime < start && last_prime != prime; {
		last_prime = prime
		src = sieveChannel(src, start, prime)
		prime = nextPrime(src, prime)
	}
	if start <= 2 {
		dst <- 2
		if end > 0 {
			dst <- prime
		}
	}
	if end == 0 {
		for range src {
		}
	}
	for last_prime := uint(0); last_prime != prime; {
		last_prime = prime
		src = sieveChannel(src, start, prime)
		prime = nextPrimeSend(src, dst, prime)
	}
}

func startSieve(start, end uint) <-chan uint {
	src := make(chan uint, 1)
	go generate(src, start, end)
	dst := make(chan uint, 1)
	go filterRoutine5(src, dst, start, end, <-src)
	return dst
}

func main() {
	var start uint
	var end uint
	flag.UintVar(&start, "start", 0, "usage")
	flag.UintVar(&end, "end", 0, "usage")
	flag.Parse()
	if end > 0 && end <= 2 {
		return
	}
	for prime := range startSieve(start, end) {
		fmt.Print(prime, "\n")
	}
}
