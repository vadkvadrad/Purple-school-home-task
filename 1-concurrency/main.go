package main

import (
	"fmt"
	"math/rand"
)

func main() {
	random := make(chan int)
	square := make(chan int)
	go Random(random)
	go Square(random, square)
	for number := range square{
		fmt.Printf("%d ", number)
	}
}

func Random(out chan<- int) {
	numbers := make([]int, 10)
	for i := range numbers {
		numbers[i] = rand.Intn(100)
	}
	for _, v := range numbers {
		out <- v
	}
	close(out)
}

func Square(in <-chan int, out chan<- int) {
	for number := range in {
		out <- number*number
	}
	close(out)
}