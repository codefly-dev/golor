# color
Turns github.com/fatih/color into HTML. Just kidding!

## Usage

```go
package main

import (
	"fmt"

	render "github.com/hygge-io/color"
)

func main() {

	s := `This is a #red{part of text with #bold{some} in bold} word.
Possible to #(green,italic){combine}`

	renderer := render.New()
	fmt.Println(renderer.Render(s))
}
```


