// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/seehuhn/password"
)

func main() {
	timings := make(chan time.Time)
	go func(timings <-chan time.Time) {
		base := <-timings
		for t := range timings {
			dt := t.Sub(base)
			fmt.Println(dt)
			base = t
		}
	}(timings)

	fmt.Println("before")
	input, err := password.ReadWithTimings("passwd: ", timings)
	fmt.Println("after")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("read %q %v\n", string(input), input)
	}
}
