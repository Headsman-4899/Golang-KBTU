package main

import (
	"fmt"
	"sync"
	"time"
)

func sumOfDigits(number int) <-chan int {
	ch := make(chan int)

	sum := 0
	n := number

	go func() {
		for n != 0 {
			digit := n % 10
			sum += digit
			n /= 10
		}
		time.Sleep(2 * time.Second)

		ch <- sum
		close(ch)
	}()

	return ch
}

func merge(channels ...<-chan int) <-chan int {
	ch := make(chan int)
	wg := new(sync.WaitGroup)

	for _, c := range channels {
		wg.Add(1)

		localC := c
		go func() {
			defer wg.Done()

			for in := range localC {
				ch <- in
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func main() {
	mergedChannels := merge(sumOfDigits(12), sumOfDigits(34))
	for value := range mergedChannels {
		fmt.Println(value)
	}

	fmt.Println("Done...")
}
