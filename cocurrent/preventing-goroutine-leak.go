package main
import "fmt"
import "time"

func main() {
	doWork := func(done <- chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("dowork done")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Canceling dowork goroutine")
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}
