package main

func main() {
	done := make(chan interface{})
	stringStream := make(chan string, 5)
	for _, s := range []string{"a", "b", "c"}{
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}

	for {
		select {
		case <-done:
			return
		default:
			//do something
		}
	}
}