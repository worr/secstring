# Secure-ish strings in Go!

secstring aims to provide a basic secure string implementation to go

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
