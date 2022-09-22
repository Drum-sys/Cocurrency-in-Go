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

### Or-channel
一般将多个已完成channel组合为一个已完成channel，当其中一个组件通道关闭，该通道就会关闭
```go
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
```

### Error handling
当goroutine当中错误太多时， 应该将错误封装保存起来发送给程序的另一部分，该部分拥有程序状态的完整信息。
func main() {
	checkStatus := func(done <- chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls{
				var result Result
				rsp, err := http.Get(url)
				//if err != nil {
				//	fmt.Println(err)
				//}
				result = Result{
					Error: err,
					Response: rsp,
				}
				select {
				case <-done:
					return
				case results <-result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...){
		if result.Error != nil {
			fmt.Printf("error : %v", result.Error)
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}


### Pipelines
使用channel构造pipeline
```go
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
```
![图片](https://user-images.githubusercontent.com/82791037/191648285-997b7e1d-16f8-4e00-b0f9-ddc86828e8a3.png)


### fanIn and fanOut
如果一个pipeline运行时间过长，或者某个stage计算代价昂贵，我们可以重用pipeline的每个stage， 并行化这个stage提升pipeline的性能。
```go
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
```

### Or-done-channel
当代码通过done管道取消时，不知道goroutine取消是否意味着正在读取的channel也被取消。
```go
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
					// 添加select语句判断 
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
```

### Tee-channel
将一个channel分离成两个channel
```go
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
```	
	




