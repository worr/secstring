package main

import (
	"fmt"
	"github.com/worr/secstring"
	"log"
)

func main() {
	str := "Man do I love safe strings!"
	ss, err := secstring.FromString(&str)
	if err != nil {
		log.Fatalf("Can't initialize string: %v", err)
	}
	defer ss.Destroy()

	fmt.Printf("Old string: %v\n", str)
	fmt.Printf("String: %v\n", string(ss.String))
}
