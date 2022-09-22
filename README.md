# Cocurrency-in-Go
## Cocurrency Patterns
### for-select loop
```go
	done := make(chan interface{})
	stringStream := make(chan string, 5)
  // 在chan上发送迭代变量
	for _, s := range []string{"a", "b", "c"}{
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
  // 循环等待退出条件
	for {
		select {
		case <-done:
			return
		default:
			//do something
		}
	}
```

### Preventing Goroutine Leaks
父goroutine应该告诉子goroutine什么时候退出，一般情况下传递一个无缓冲的channel，结合for select使用，当关闭无缓冲channel时，select接收到关闭信号，return。避免goroutine一直留在内存中
```go
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
```

###
