// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/seehuhn/password"
)

func main() {
	fmt.Println("before")
	input, err := password.Read("passwd: ")
	fmt.Println("after")
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("read %q %v\n", string(input), input)
	}
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("now!")
	}()
}
