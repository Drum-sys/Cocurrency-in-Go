package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, inters ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range inters{
				select {
				case <-done:
					return
				case intStream <-i:
				}
			}
		}()
		return intStream
	}

	mul := func(done <-chan interface{}, intstream <-chan int, mulVal int) <-chan int{
		mulStream := make(chan int)
		go func() {
			defer close(mulStream)
			for i := range intstream {
				select {
				case <-done:
					return
				case mulStream <-i * mulVal:
				}
			}
		}()
		return mulStream
	}

	add := func(done <-chan interface{}, intstream <-chan int, addVal int) <-chan int{
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for i := range intstream {
				select {
				case <-done:
					return
				case addStream <- i + addVal:
				}
			}
		}()
		return addStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := mul(done, add(done, mul(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}

}