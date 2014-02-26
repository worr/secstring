# Secure-ish strings in Go!

secstring aims to provide a basic secure string implementation to go

## Badges!

[![Build Status](https://travis-ci.org/worr/secstring.png?branch=master)](https://travis-ci.org/worr/secstring)
[![GoDoc](https://godoc.org/github.com/worr/secstring?status.png)](https://godoc.org/github.com/worr/secstring)
[![Flattr Button](http://api.flattr.com/button/button-compact-static-100x17.png "Flattr This!")](https://flattr.com/submit/auto?user_id=worr&url=https%3A%2F%2Fgithub.com%2Fworr%2Fsecstring%2F "secstring")

## What makes them secure?

* strings are unlikely to be written to swap (except during hibernation)
* strings are immutable - modifying them causes a non-recoverable `panic`
* strings are encrypted in memory

## Wait, isn't the key in memory too?

Yes. I'm not promising perfect security. Mostly I aim to prevent trivially
grabbing the string from memory, or modifying it while it's in memory.

## This doesn't work on OSX/Windows/FreeBSD/etc.

Yes. I use `syscall` heavily, and unfortunately, golang in osx and Windows
don't have the functions I'm using. I'm going to submit patches, so hopefully
they get added soon.

## Can I get an example?

Damn straight.

```go
import "github.com/worr/secstring"
import "fmt"

func main() {
    str := "testing"
    ss, _ := secstring.FromString(&str)
    defer ss.Destroy()

    ss.Decrypt()
    fmt.Printf("String: %v", ss.String)
    ss.Encrypt()
}
```
