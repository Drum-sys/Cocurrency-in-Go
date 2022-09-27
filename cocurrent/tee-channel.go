package main

import "fmt"

func main() {

	orDone := func(done <-chan interface{}, c <-chan int) <-chan int {
		valStream := make(chan int)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					fmt.Println("done4")
					return
				case v, ok := <-c:
					if ok == false {
						fmt.Println("已经结束了")
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
	//
	//repeat := func(
	//	done <-chan interface{},
	//	values ...int,
	//) <-chan int {
	//	valueStream := make(chan int)
	//	go func() {
	//		defer close(valueStream)
	//		for {
	//			for _, v := range values {
	//				select {
	//				case <-done:
	//					fmt.Println("done1")
	//					return
	//				case valueStream <- v:
	//				}
	//			}
	//		}
	//	}()
	//	return valueStream
	//}
	//
	//take := func(
	//	done <-chan interface{},
	//	valueStream <-chan int,
	//	num int,
	//) <-chan int {
	//	takeStream := make(chan int)
	//	go func() {
	//		defer close(takeStream)
	//		for i := 0; i < num; i++ {
	//			select {
	//			case <-done:
	//				fmt.Println("done2")
	//				return
	//			case takeStream <- <- valueStream:
	//			}
	//		}
	//	}()
	//	return takeStream
	//}

	tee := func(done <-chan interface{}, in <-chan int) (<-chan int, <-chan int) {
		out1 := make(chan int)
		out2 := make(chan int)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				//fmt.Println(val)
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						fmt.Println("done3")
						//return
					case out1 <-val:
						out1 = nil
					case out2 <-val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	mychan := make(chan int, 5)
	mychan <- 1
	mychan <- 2
	mychan <- 3
	done := make(chan interface{})
	defer close(done)
	//out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))
	out1, out2 := tee(done, mychan)
	close(mychan)
	for val := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val, <-out2)
	}


}
