package main

import (
	"golox/parser"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		wantType string
		wantLen  int
	}{
		{
			name:    "simple number",
			source:  "123",
			wantLen: 2, // NUMBER + EOF
		},
		{
			name:    "simple string",
			source:  `"hello"`,
			wantLen: 2, // STRING + EOF
		},
		{
			name:    "boolean keywords",
			source:  "true false",
			wantLen: 3, // TRUE FALSE EOF
		},
		{
			name:    "arithmetic expression",
			source:  "1 + 2 * 3",
			wantLen: 6, // NUMBER PLUS NUMBER STAR NUMBER EOF
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.source)
			tokens := scanner.ScanTokens()

			if len(tokens) != tt.wantLen {
				t.Errorf("got %d tokens, want %d", len(tokens), tt.wantLen)
			}
		})
	}
}

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "simple arithmetic",
			source: "1 + 2",
			want:   "3",
		},
		{
			name:   "multiplication precedence",
			source: "2 + 3 * 4",
			want:   "14",
		},
		{
			name:   "grouping",
			source: "(2 + 3) * 4",
			want:   "20",
		},
		{
			name:   "boolean equality",
			source: "true == true",
			want:   "true",
		},
		{
			name:   "comparison",
			source: "5 > 3",
			want:   "true",
		},
		{
			name:   "string concatenation",
			source: `"Hello" + " " + "World"`,
			want:   "Hello World",
		},
		{
			name:   "nil",
			source: "nil",
			want:   "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// You'll need to capture output for proper testing
			// For now, just check it doesn't panic
			scanner := NewScanner(tt.source)
			tokens := scanner.ScanTokens()
			parser := parser.NewParser(tokens)
			expr, err := parser.Parse()

			if err != nil {
				t.Errorf("parse error: %v", err)
				return
			}

			if expr == nil {
				t.Error("got nil expression")
				return
			}

			// Would need to capture interpreter output to test properly
			interpreter := NewInterpreter()
			interpreter.Interpret(expr)
		})
	}
}
