package render_test

import (
	"fmt"
	render "github.com/hygge-io/color"
	"testing"
)

func TestScanner(t *testing.T) {
	scanner := render.NewScanner()
	tcs := []struct {
		name string
		text string
		want []render.Token
	}{
		{name: "no rendering", text: "no rendering", want: []render.Token{{Text: "no rendering"}}},
		{name: "simple", text: "this #red{cat is colorful} and this is not", want: []render.Token{{Text: "this "}, {Text: "cat is colorful", Style: render.NewStyle().Color(render.Red)}, {Text: " and this is not"}}},
		{name: "two tags", text: "this #(red,bold){cat is colorful and awesome} and this is not", want: []render.Token{{Text: "this "}, {Text: "cat is colorful and awesome", Style: render.NewStyle().Color(render.Red).Typography(render.Bold)}, {Text: " and this is not"}}},
		{name: "nested", text: "this #red{cat is #bold{awesome}}", want: []render.Token{{Text: "this "}, {Text: "cat is ", Style: render.NewStyle().Color(render.Red)}, {Text: "awesome", Style: render.NewStyle().Color(render.Red).Typography(render.Bold)}}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := scanner.Scan(tc.text)
			if len(got) != len(tc.want) {
				t.Errorf("got %d tokens, want %d", len(got), len(tc.want))
			}
			for i, g := range got {
				if g.Text != tc.want[i].Text {
					t.Errorf("Text not matching, got <%v>, want <%v>", g.Text, tc.want[i].Text)
				}
				if !render.SameStyle(g.Style, tc.want[i].Style) {
					t.Errorf("Style not matching for token <%s>, got %v, want %v", g.Text, g.Style, tc.want[i].Style)
				}
			}
		})
	}
	s := `This is a #red{part of text with #bold{some} in bold} word.
Possible to #(blue,italic){combine}`

	renderer := render.New()
	fmt.Println(renderer.Render(s))
}
