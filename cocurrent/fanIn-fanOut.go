package main
import (
	"sync"
)

func main() {

	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		mulStream := make(chan interface{})
		mulPlex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case mulStream <-i:
				}
			}
		}
		wg.Add(len(channels))
		for _, c := range channels {
			go mulPlex(c)
		}

		go func() {
			wg.Wait()
			close(mulStream)
		}()
		return mulStream
	}
	done := make(chan interface{})
	defer close(done)
	// 自己定义channels
	<-fanIn(done, nil)
}
