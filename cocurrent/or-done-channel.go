package main

import "fmt"

func main() {
	orDone := func(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <-v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	mychan := make(chan interface{}, 4)
	done := make(chan interface{})
	defer close(done)
	for val := range orDone(done, mychan) {
		fmt.Println(val)
	}
}