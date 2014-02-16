# Jumpy

Jumpy is a versatile web cralwer library. On the other hand run stand-alone in terminal as command.

## Install

To use as library:

```
go get github.com/yoppi/jumpy
```

Otherwise, to use as command:

```
go get github.com/yoppi/jumpy/jumpy
```

## Usage

Use as library in your code.

```go
import (
  "github.com/yoppi/jumpy"
  "fmt"
)

func main() {
  jumpy.Crawl("http://golang.org/", map[string]string{}, func(page *jumpy.Page) {
    fmt.Println(page.Url)
  })
}
```

Use as command in terminal.

```
$ jumpy --root=http://golang.org/ --command=url
```

