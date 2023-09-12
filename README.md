# color
Turns github.com/fatih/color into HTML. Just kidding!

Just a tag based approach to write some nice text.

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

![output](media/output.png)

Note:

The code enables to change the triplet `# { }` to define tagging. Just need to expose it in the constructor.
