# wallint
A minimum golint example to find variables named wally.

```
$ go get -u https://github.com/minamijoyo/wallint
```

```go
package main

import "fmt"

func main() {
	wally := 1
	fmt.Println(wally)
}
```

```
$ wallint data/wally.go
data/wally.go:6:2: wally was found
```
