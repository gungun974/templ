package parser

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestTextParser(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected Text
	}{
		{
			name:  "Text ends at an element start",
			input: `abcdef<a href="https://example.com">More</a>`,
			expected: Text{
				Value:     "abcdef",
				LineBreak: true,
			},
		},
		{
			name:  "Text ends at a templ expression start",
			input: `abcdef{%= "test" %}`,
			expected: Text{
				Value:     "abcdef",
				LineBreak: true,
			},
		},
		{
			name:  "Text may contain spaces",
			input: `abcdef ghijk{%= "test" %}`,
			expected: Text{
				Value:     "abcdef ghijk",
				LineBreak: true,
			},
		},
		{
			name:  "Text may contain named references",
			input: `abcdef&nbsp;ghijk{%= "test" %}`,
			expected: Text{
				Value:     "abcdef&nbsp;ghijk",
				LineBreak: true,
			},
		},
		{
			name:  "Text may contain base 10 numeric references",
			input: `abcdef&#32;ghijk{%= "test" %}`,
			expected: Text{
				Value:     "abcdef&#32;ghijk",
				LineBreak: true,
			},
		},
		{
			name:  "Text may contain hexadecimal numeric references",
			input: `abcdef&#x20;ghijk{%= "test" %}`,
			expected: Text{
				Value:     "abcdef&#x20;ghijk",
				LineBreak: true,
			},
		},
		{
			name:  "Text should not LineBreak when it's next to an opening none line breaking element",
			input: `abcd</span>the_end`,
			expected: Text{
				Value:     "abcd",
				LineBreak: false,
			},
		},
		{
			name:  "Text should not LineBreak when it's next to an closing none line breaking element",
			input: `abcd<span>the_end</span>`,
			expected: Text{
				Value:     "abcd",
				LineBreak: false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			actual, ok, err := textParser.Parse(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("unexpected failure for input %q", tt.input)
			}
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}
