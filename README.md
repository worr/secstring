# Secure-ish strings in Go!

secstring aims to provide a basic secure string implementation to go

## Badges!

[![Build Status](https://travis-ci.org/worr/secstring.png?branch=master)](https://travis-ci.org/worr/secstring)
[![Coverage Status](https://coveralls.io/repos/worr/secstring/badge.png)](https://coveralls.io/r/worr/secstring)
[![GoDoc](https://godoc.org/gitlab.com/worr/secstring.git?status.png)](https://godoc.org/github.com/worr/secstring)

## Should I use this?

Probably not. I've implemented this mostly as a PoC. I use it somewhat, but I don't recommend other people use it right now.

## What makes them secure?

* strings are unlikely to be written to swap (except during hibernation)
* strings are immutable - modifying them causes a non-recoverable `panic`

## This doesn't work on Windows/FreeBSD/etc.

Yes. I use `syscall` heavily, and unfortunately, golang in many BSDs
don't have the functions I'm using. I'm going to submit patches, so hopefully
they get added soon.

Windows support will never be added. I don't have a test box for it.

## Can I get an example?

Damn straight.

```go
import "github.com/worr/secstring"
import "fmt"

func main() {
    str := "testing"
    ss, _ := secstring.FromString(&str)
    defer ss.Destroy()

    fmt.Printf("String: %v", ss.String)
}
```
