package main

import (
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hex conversion",
			input:    "1E (hex) files were added",
			expected: "30 files were added",
		},
		{
			name:     "bin conversion",
			input:    "It has been 10 (bin) years",
			expected: "It has been 2 years",
		},
		{
			name:     "up case",
			input:    "Ready, set, go (up) !",
			expected: "Ready, set, GO!",
		},
		{
			name:     "low case",
			input:    "I should stop SHOUTING (low)",
			expected: "I should stop shouting",
		},
		{
			name:     "cap case",
			input:    "Welcome to the Brooklyn bridge (cap)",
			expected: "Welcome to the Brooklyn Bridge",
		},
		{
			name:     "numbered up",
			input:    "This is so exciting (up, 2)",
			expected: "This is SO EXCITING",
		},
		{
			name:     "numbered cap",
			input:    "it was the age of wisdom, it was the age of foolishness (cap, 6)",
			expected: "it was the age of wisdom, It Was The Age Of Foolishness",
		},
		{
			name:     "punctuation spacing",
			input:    "I was sitting over there ,and then BAMM !!",
			expected: "I was sitting over there, and then BAMM!!",
		},
		{
			name:     "punctuation groups",
			input:    "I was thinking ... You were right",
			expected: "I was thinking... You were right",
		},
		{
			name:     "single quotes",
			input:    "I am exactly how they describe me: ' awesome '",
			expected: "I am exactly how they describe me: 'awesome'",
		},
		{
			name:     "multi word quotes",
			input:    "As Elton John said: ' I am the most well-known homosexual in the world '",
			expected: "As Elton John said: 'I am the most well-known homosexual in the world'",
		},
		{
			name:     "article correction",
			input:    "There it was. A amazing rock!",
			expected: "There it was. An amazing rock!",
		},
		{
			name:     "mixed hex up",
			input:    "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			expected: "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			name:     "full sample",
			input:    "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			expected: "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			name:     "punctuation tests",
			input:    "Punctuation tests are ... kinda boring ,what do you think ?",
			expected: "Punctuation tests are... kinda boring, what do you think?",
		},
		{
			name:     "empty line",
			input:    "line1\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "leading modifier",
			input:    "(up) first word",
			expected: "first word",
		},
		{
			name:     "invalid number",
			input:    "abg (hex)",
			expected: "abg",
		},
		{
			name:     "valid hex letters",
			input:    "abc (hex)",
			expected: "2748",
		},
		{
			name:     "large numbered",
			input:    "short (up, 10)",
			expected: "SHORT",
		},
		{
			name:     "article with h",
			input:    "A hour",
			expected: "An hour",
		},
		{
			name:     "article non vowel",
			input:    "a banana",
			expected: "a banana",
		},
		{
			name:     "quote with punct",
			input:    "say 'hello' !",
			expected: "say 'hello'!",
		},
		{
			name:     "multi punct !?.",
			input:    "wow !?.",
			expected: "wow!?.",
		},
		{
			name:     "article in quotes with bin and punct",
			input:    "it was 10 (bin) 'a Apple'!?",
			expected: "it was 2 'an Apple'!?",
		},
		{
			name:     "article non-vowel in quotes",
			input:    "it was 'a banana'!",
			expected: "it was 'a banana'!",
		},
		{
			name:     "cap article vowel",
			input:    "It was A Apple",
			expected: "It was An Apple",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processText(tt.input)
			if got != tt.expected {
				t.Errorf("\ninput:\n%s\n\ngot:\n%s\n\nwant:\n%s", tt.input, got, tt.expected)
			}
		})
	}
}
