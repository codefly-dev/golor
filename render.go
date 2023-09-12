package render

import (
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"strings"
)

type Renderer struct {
	theme   Theme
	scanner *Scanner
}

func (renderer *Renderer) Render(text string) string {
	tokens := renderer.scanner.Scan(text)
	var rendered []string
	for _, token := range tokens {
		rendered = append(rendered, renderer.theme.Convert(token))
	}
	return strings.Join(rendered, "")
}

func New() *Renderer {
	return &Renderer{
		theme:   Theme{},
		scanner: NewScanner(),
	}
}

type Scanner struct {
	TagMarker int32
	Start     int32
	End       int32
}

func NewScanner() *Scanner {
	return &Scanner{TagMarker: '#', Start: '{', End: '}'}
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
	color        Color
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
	if a == nil && b != nil {
		return false
	}
	if a.color != b.color {
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
	style.color = Color{value: color}
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

type Theme struct {
	// Could be interesting to have primary/secondary...and theme
	// Dracula...
}

func (theme *Theme) Color(c Color) (color.Attribute, bool) {
	switch c.value {
	case Black:
		return color.FgBlack, true
	case Red:
		return color.FgRed, true
	case Green:
		return color.FgGreen, true
	case Yellow:
		return color.FgYellow, true
	case Blue:
		return color.FgBlue, true
	case Magenta:
		return color.FgMagenta, true
	case Cyan:
		return color.FgCyan, true
	case White:
		return color.FgWhite, true
	default:
		return 0, false // Return false for unsupported color indices
	}
}

func (theme *Theme) Typographies(typographies []Typography) []color.Attribute {
	var attributes []color.Attribute
	for _, t := range typographies {
		switch t.value {
		case Bold:
			attributes = append(attributes, color.Bold)
		case Italic:
			attributes = append(attributes, color.Italic)
		}
	}
	return attributes
}

func (theme *Theme) Convert(token Token) string {
	// Use fatih/color
	if token.Style == nil {
		return token.Text
	}
	var attributes []color.Attribute
	if c, ok := theme.Color(token.Style.color); ok {
		attributes = append(attributes, c)
	}
	attributes = append(attributes, theme.Typographies(token.Style.typographies)...)
	f := color.New(attributes...).SprintFunc()
	return f(token.Text)
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
	if level != 0 {
		panic(fmt.Sprintf("invalid number of tags: %d", level))
	}

	// Add the last token
	if len(currentText) > 0 {
		tokens = append(tokens, Token{Text: currentText, Style: styles[level]})
	}
	return tokens
}
