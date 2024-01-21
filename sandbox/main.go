package main

import (
	"github.com/codefly-dev/golor"
)

type display struct {
	FromTemplate string
	Items        []string
}

func Manual() {
	golor.Println(`This is a #red[part of text with #bold[some] in bold] word`)
	golor.Println(`This is a #(blue)[other %s] with regular %s`, "text", "arguments")

	golor.Println(`This is a #red[part of text with #bold[some] in bold] word. Possible to #(blue,italic)[combine]`)

	// Work with a template
	golor.Template(display{FromTemplate: "Hello from template",
		Items: []string{"Item 1", "Item 2"},
	}).Println(`
#(blue)[Welcome to #(red)[codefly-io/golor]!]
#(green,italic)[{{.FromTemplate}} in italic]
#(cyan,bold)[<some brackets> in bold]
{{- range .Items}}
#(yellow,bold)[{{.}}]{{end}}
Notice the #(white,bold,italic)[trick] to not have lines between the range items!`)

	s := `
In Markdown, @bold<a new paragraph> uses the # tag
while links are written as @green<[link](url)>`
	renderer := golor.New().WithTagMarker('@').WithTextLimiter('<', '>')
	renderer.Println(s)
}

func main() {
	Manual()
}
