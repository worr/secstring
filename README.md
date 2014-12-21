# Secure-ish strings in Go!

secstring aims to provide a basic secure string implementation to go

## Badges!

[![Build Status](https://travis-ci.org/worr/secstring.png?branch=master)](https://travis-ci.org/worr/secstring)
[![Coverage Status](https://coveralls.io/repos/worr/secstring/badge.png)](https://coveralls.io/r/worr/secstring)
[![GoDoc](https://godoc.org/gitlab.com/worr/secstring.git?status.png)](https://godoc.org/github.com/worr/secstring)
[![Flattr Button](http://api.flattr.com/button/button-compact-static-100x17.png "Flattr This!")](https://flattr.com/submit/auto?user_id=worr&url=https%3A%2F%2Fgithub.com%2Fworr%2Fsecstring%2F "secstring")

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
import "gitlab.com/worr/secstring.git"
import "fmt"

func main() {
    str := "testing"
    ss, _ := secstring.FromString(&str)
    defer ss.Destroy()

    fmt.Printf("String: %v", ss.String)
}
```
