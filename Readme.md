# Useful Colors [![Go Reference](https://pkg.go.dev/badge/grow.graphics/uc.svg)](https://pkg.go.dev/grow.graphics/uc)

This module provides a useful purego color representation for graphics programming with zero dependencies, 
the design and implementation follows Godot's [Color Class](https://docs.godotengine.org/en/stable/classes/class_color.html).

```go
package main

import (
    "fmt"

    "grow.graphics/uc"
)

func main() {
    fmt.Println(uc.RGB(1, 0, 0))
    fmt.Println(uc.X11.Violet)
    fmt.Println(uc.Hex(0xff0000))
    fmt.Println(uc.HTML("#ff0000"))
}
```