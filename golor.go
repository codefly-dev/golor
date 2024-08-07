package golor

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

// Default Renderer helpers

func Println(s string, args ...any) {
	renderer := New()
	txt := fmt.Sprintf(s, args...)
	fmt.Println(renderer.Render(txt))
}

func Sprintf(s string, args ...any) string {
	renderer := New()
	txt := fmt.Sprintf(s, args...)
	return renderer.Render(txt)
}

func Template(t any) *Renderer {
	return New().WthTemplate(t)
}

type Renderer struct {
	theme    Theme
	scanner  *Scanner
	template any
}

func (renderer *Renderer) Render(text string) string {
	var tokens []Token
	if renderer.template == nil {
		tokens = renderer.scanner.Scan(text)
	} else {
		tokens = renderer.scanner.ScanWithTemplate(text, renderer.template)
	}
	return renderer.theme.Produce(tokens)
}

func New() *Renderer {
	return &Renderer{
		theme:   Theme{},
		scanner: NewScanner(),
	}
}

func (renderer *Renderer) Println(s string, args ...any) {
	fmt.Println(renderer.Render(fmt.Sprintf(s, args...)))
}

func (renderer *Renderer) Sprintf(s string, args ...any) string {
	return renderer.Render(fmt.Sprintf(s, args...))
}

func (renderer *Renderer) Scanner() *Scanner {
	return renderer.scanner
}

func (renderer *Renderer) WithTagMarker(m int32) *Renderer {
	renderer.scanner.TagMarker = m
	return renderer
}

func (renderer *Renderer) WithTextLimiter(start, end int32) *Renderer {
	renderer.scanner.Start = start
	renderer.scanner.End = end
	return renderer
}

func (renderer *Renderer) WthTemplate(t any) *Renderer {
	renderer.template = t
	return renderer
}

type Scanner struct {
	TagMarker int32
	Start     int32
	End       int32
}

func NewScanner() *Scanner {
	return &Scanner{TagMarker: '#', Start: '[', End: ']'}
}

type Token struct {
	Text  string
	Style *Style
}

type Type interface {
}

const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Color struct {
	value int
}

const (
	None = iota
	Bold
	Italic
)

type Typography struct {
	value int
}

type Style struct {
	color        *Color
	typographies []Typography
}

func (style *Style) Clone() *Style {
	s := Style{
		color: style.color,
	}
	copy(s.typographies, style.typographies)
	return &s
}

func (style *Style) Add(ts ...string) {
	for _, t := range ts {
		style.AddOne(t)
	}
}

func (style *Style) String() string {
	return fmt.Sprintf("color: %d, typographies: %v", style.color, style.typographies)
}

func NewStyle() *Style {
	return &Style{}
}

func SameStyle(a, b *Style) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil {
		return false
	}
	if !reflect.DeepEqual(a.color, b.color) {
		return false
	}
	if len(a.typographies) != len(b.typographies) {
		return false
	}
	match := make(map[int]bool)
	for _, t := range a.typographies {
		match[t.value] = true
	}
	for _, t := range b.typographies {
		if !match[t.value] {
			return false
		}
	}
	return true
}

func (style *Style) Color(color int) *Style {
	style.color = &Color{value: color}
	return style
}

func (style *Style) Typography(t int) *Style {
	existing := make(map[int]bool)
	for _, t := range style.typographies {
		existing[t.value] = true
	}
	if !existing[t] {
		style.typographies = append(style.typographies, Typography{value: t})
	}
	return style
}

var colorNames = map[string]int{
	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
}

func (style *Style) AddOne(t string) {
	switch t {
	case "bold":
		style.Typography(Bold)
	case "italic":
		style.Typography(Italic)
	default:
		if c, ok := colorNames[t]; ok {
			style.Color(c)
		}
	}
}

func (s *Scanner) Scan(text string) []Token {
	var tokens []Token
	var currentText string
	var currentTags string
	var insideTag bool

	styles := make(map[int]*Style)
	level := 0
	reg := regexp.MustCompile(`\(?([\w,]*)\)?`)

	for _, char := range text {
		switch char {
		case s.TagMarker:
			token := Token{Text: currentText, Style: styles[level]}
			tokens = append(tokens, token)
			insideTag = true
			currentText = ""
		case s.Start:
			match := reg.Match([]byte(currentTags))
			if !match {
				panic(fmt.Sprintf("invalid tag: %s", currentTags))
			}
			currentTags = reg.FindStringSubmatch(currentTags)[1]
			tags := strings.Split(currentTags, ",")
			previous := styles[level]
			level += 1
			next := NewStyle()
			if previous != nil {
				next = previous.Clone()
			}
			next.Add(tags...)
			styles[level] = next
			currentTags = ""
			insideTag = false
		case s.End:
			if len(currentText) > 0 {
				token := Token{Text: currentText, Style: styles[level]}
				tokens = append(tokens, token)
			}
			currentText = ""
			level -= 1

		default:
			if insideTag {
				currentTags += string(char)
			} else {
				currentText += string(char)
			}
		}
	}

	// Add the last token
	if len(currentText) > 0 {
		tokens = append(tokens, Token{Text: currentText, Style: styles[level]})
	}
	return tokens
}

func (s *Scanner) ScanWithTemplate(text string, obj any) []Token {
	tmpl, err := template.New("rendering").Parse(text)
	if err != nil {
		panic(fmt.Sprintf("cannot parse template: %s", err))
	}
	var b strings.Builder
	err = tmpl.Execute(&b, obj)
	if err != nil {
		panic(fmt.Sprintf("cannot execute template: %s", err))
	}
	return s.Scan(b.String())
}
