package main

import (
	"fmt"
	"github.com/hygge-io/golor"
)

func main() {

	golor.Println(`This is a #red[part of text with #bold[some] in bold] word.
Possible to #(blue,italic)[combine]`)

	golor.Println(`#(blue)[Welcome to #(red)[hygge-io/golor]!]
#(green,italic)[{{.FromTemplate}} in italic]
#(cyan,bold)[<some brackets> in bold]
{{- range .Items}}
#(yellow,bold)[{{.}}]{{end}}
Notice the #(white,bold,italic)[trick] to not have lines between the range items!`, map[string]any{
		"FromTemplate": "Hello from template",
		"Items":        []string{"Item 1", "Item 2"},
	})

	s := `In Markdown, @(bold)<a new paragraph> uses the # tag
while links are written as @green<[link](url)>`
	renderer := golor.New().WithTagMarker('@').WithTextLimiter('<', '>')
	fmt.Println(renderer.Render(s))
}
