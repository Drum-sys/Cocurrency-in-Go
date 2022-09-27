package main

import (
	"net/http"
	"fmt"
)

type Result struct {
	Error error
	Response *http.Response
}

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
