package main

import "time"
import "fmt"

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <- chan interface{} {
		fmt.Printf("channel length %v\n", len(channels))
		switch len(channels) {
		case 0:
			//fmt.Println("第一次静茹case0")
			return nil
		case 1:
			//fmt.Println("第一次静茹case0")
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				fmt.Println("长度等于2")
				select {
				case <-channels[0]:
					fmt.Println("第一次静茹case0")
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
					fmt.Println("第一次静茹case0")
				case <-channels[1]:
					fmt.Println("第一次静茹case1")
				case <-channels[2]:
					fmt.Println("第一次静茹case2")
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		fmt.Println(">>>>>>>>>>>>>>>>>>")
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(1*time.Second),
		sig(1*time.Second),
		//sig(2*time.Hour),
		sig(5*time.Millisecond),
		sig(1*time.Second),
		//sig(1*time.Hour),
		//sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))

	time.Sleep(3 * time.Second)

	arr := []int{1,2}
	fmt.Println(append(arr, 3))

}
