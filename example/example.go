package main

import "gitlab.com/worr/secstring.git"
import "fmt"

func main() {
    str := "Man do I love obfuscated strings!"
    ss, _ := secstring.FromString(&str)
    defer ss.Destroy()

    ss.Decrypt()
    fmt.Printf("String: %v\n", string(ss.String))
    ss.Encrypt()
    fmt.Printf("String: %v\n", string(ss.String))
    ss.Decrypt()
    fmt.Printf("String: %v\n", string(ss.String))
    ss.Encrypt()
    fmt.Printf("String: %v\n", string(ss.String))
}
