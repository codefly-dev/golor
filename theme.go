package golor

import (
	"github.com/fatih/color"
	"strings"
)

type Theme struct {
	// Could be interesting to have primary/secondary...and theme
	// Dracula...
}

func (theme *Theme) Produce(tokens []Token) string {
	var rendered []string
	for _, token := range tokens {
		rendered = append(rendered, theme.Convert(token))
	}
	return strings.Join(rendered, "")
}

func (theme *Theme) Color(c *Color) (color.Attribute, bool) {
	if c == nil {
		return 0, false
	}
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
