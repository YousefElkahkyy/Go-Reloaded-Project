package main

import (
	"os"
	"strconv"
	"strings"
)

/*
Milestone 1 — Read Input, Write Output
How many command-line arguments does your program expect? What happens if the user provides the wrong number?
3 arguments ([0] program name, [1] input file, [2] output file). If wrong number, print usage and exit.

What type does os.ReadFile return, and what do you need to do before you can work with it as text?
[]byte and error. Convert []byte to string(content).

What does the third argument of os.WriteFile control?
File permissions/mode (0644: owner rw, others r).

Milestone 2 — Number Conversions: (hex) and (bin)
What function converts a string like "42" from base 16 to decimal? strconv.ParseInt(s, 16, 64)
Base for hex: 16, bin: 2.
If not valid number (err != nil), return original s unchanged.

Milestone 3 — Case Modifiers: (up), (low), (cap)
Detect modifier token, transform previous word in result slice, skip modifier.
Capitalize: first upper, rest lower.

Milestone 4 — Numbered Modifiers: (up, 2), (cap, 4)
Parse (type,N): split ,, Atoi N after trim ) , apply to last min(N, len(result)) words.
Guard start >=0.

Milestone 5 — Punctuation Spacing
Punct token (. , ! ? : ; groups ... !? ): attach to previous word no space.

Milestone 6 — Single Quote Formatting
Standalone ' : opening attaches to next word, closing to last word.

Milestone 7 — Article Correction: a → an
If "a"/"A", next word starts vowel/h (lower), replace with "an"/"An"; reverse if "an"/"An" non-vowel.

Milestone 8 — Connect the Pipeline
Order: tokenize -> modifiers -> quotes -> punct -> articles -> join " ".
strings.Fields equivalent via custom tokenizer for punct/parens.
*/

func main() {
	// 1. Check that exactly 2 arguments were provided (input file and output file paths)
	//    Total len(os.Args) == 3 including program name
	//    If wrong number, print usage message and exit early
	if len(os.Args) != 3 {
		// Print usage message and return
		println("Usage: go run . <input_file> <output_file>")
		return
	}

	// os.Args[1], [2]
	inFilename := os.Args[1]
	outFilename := os.Args[2]

	// 2. Read the input file into a string
	content, err := readFile(inFilename)
	if err != nil {
		println("Error reading file:", err.Error())
		return
	}

	// 3. Transformations here (pipeline)
	result := processText(content)

	// 4. Write the result string to the output file
	err = writeFile(outFilename, result)
	if err != nil {
		println("Error writing file:", err.Error())
		return
	}

	// Success message as per guide
	println("File processing completed successfully.")
}

func readFile(filename string) (string, error) {
	// Read the entire file into memory as []byte using os.ReadFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err // Propagate error, return empty string
	}
	// Convert []byte to string and remove any \r characters to normalize line endings
	return strings.ReplaceAll(string(data), "\r", ""), nil
}

func writeFile(filename string, content string) error {
	// Write the content bytes to file with 0644 permissions (owner rw, group/world r)
	return os.WriteFile(filename, []byte(content), 0644)
}

// Milestone 2 helper: hex word before (hex) -> decimal
func hexToDecimal(s string) string {
	// Parse input string s as base-16 (hex) integer into int64
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return s // If parsing fails (not valid hex), return original unchanged
	}
	// Convert parsed int64 back to base-10 (decimal) string
	return strconv.FormatInt(n, 10)
}

// Milestone 2 helper: bin word before (bin) -> decimal
func binToDecimal(s string) string {
	// Parse input string s as base-2 (binary) integer into int64
	n, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s // If parsing fails (not valid binary), return original unchanged
	}
	// Convert parsed int64 back to base-10 (decimal) string
	return strconv.FormatInt(n, 10)
}

func processModifiers(words []string) []string {
	// Build new slice `result` with transformed words (modifiers applied, modifiers themselves skipped)
	var result []string

	for i := 0; i < len(words); i++ {
		token := words[i]

		// Check if current token is (hex): replace LAST word in result with its decimal equivalent
		if token == "(hex)" && len(result) > 0 {
			result[len(result)-1] = hexToDecimal(result[len(result)-1])
			continue // Skip adding modifier
		}

		// Check if current token is (bin): replace LAST word in result with its decimal equivalent
		if token == "(bin)" && len(result) > 0 {
			result[len(result)-1] = binToDecimal(result[len(result)-1])
			continue // Skip adding modifier
		}

		// Simple case modifiers: (up), (low), (cap) - apply to LAST word, skip modifier
		if token == "(up)" || token == "(low)" || token == "(cap)" {
			if len(result) > 0 {
				result[len(result)-1] = applyCase(result[len(result)-1], token)
			}
			continue
		}

		// Numbered modifiers like (up, 2): parse and apply to LAST N words
		if strings.HasPrefix(token, "(") && strings.Contains(token, ",") {
			// Clean token: remove spaces, trim parentheses
			clean := strings.Trim(strings.ReplaceAll(token, " ", ""), "()")
			parts := strings.Split(clean, ",")
			if len(parts) == 2 {
				// Parse count N from second part (trim extra chars)
				n, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(parts[1]), ")"))
				if err == nil {
					// Calculate start index for last N words (guard against negative)
					start := len(result) - n
					if start < 0 {
						start = 0
					}
					modType := "(" + strings.TrimSpace(parts[0]) + ")"
					// Apply case change to last N words
					for j := start; j < len(result); j++ {
						result[j] = applyCase(result[j], modType)
					}
				}
			}
			continue // Skip adding modifier
		}

		// Normal word or token: append unchanged to result
		result = append(result, token)
	}

	return result
}

// Helper: Apply case transformation based on modifier type
// (up)=UPPER, (low)=lower, (cap)=Capitalize first letter
func applyCase(s, mod string) string {
	if s == "" {
		return s // Nothing to change
	}
	switch mod {
	case "(up)":
		return strings.ToUpper(s) // Convert entire string to UPPERCASE
	case "(low)":
		return strings.ToLower(s) // Convert entire string to lowercase
	case "(cap)":
		lower := strings.ToLower(s)
		if len(lower) == 0 {
			return s
		}
		// Capitalize: First char UPPER + rest lower
		return strings.ToUpper(string(lower[0])) + lower[1:]
	}
	return s // No matching mod, return unchanged
}

func fixPunctuation(words []string) []string {
	// Attach punctuation tokens (. , ! ? : ; and groups like !? ...) to previous word without space
	var result []string
	punctuationChars := ".,!?:;" // Defined set of punctuation starters

	for _, word := range words {
		if word == "" {
			continue // Skip empty
		}
		if len(word) > 0 && strings.ContainsRune(punctuationChars, rune(word[0])) {
			// Current word is punctuation or punct group: attach to last word in result
			if len(result) > 0 {
				result[len(result)-1] += word // Append directly (no space)
			} else {
				// No previous word: add punct as standalone
				result = append(result, word)
			}
		} else {
			// Normal word: append normally
			result = append(result, word)
		}
	}
	return result
}

func fixQuotes(words []string) []string {
	// Process standalone single quotes: opening ' attaches to NEXT word start
	// closing ' attaches to END of last word. No extra spaces around quotes.
	wordsCopy := make([]string, len(words)) // Working copy (avoid mutate input)
	copy(wordsCopy, words)
	var result []string
	quoteOpen := false // Track if inside quoted section (between opening/closing ')

	for i := 0; i < len(wordsCopy); i++ {
		if wordsCopy[i] == "'" { // Found standalone quote token
			if !quoteOpen {
				// Opening quote: prefix to next token (if exists)
				if i+1 < len(wordsCopy) {
					wordsCopy[i+1] = "'" + wordsCopy[i+1]
				}
				quoteOpen = true
			} else {
				// Closing quote: suffix to most recent word in result
				if len(result) > 0 {
					result[len(result)-1] += "'"
				}
				quoteOpen = false
			}
			continue // Skip adding the bare ' token itself
		}
		// Non-quote token: append (may have been prefixed if opening quote)
		result = append(result, wordsCopy[i])
	}
	return result
}

func fixArticles(words []string) []string {
	// Correct "a/An" ↔ "an" based on next word starting with vowel (aeiouh) or not
	// Ignores attached punct/quotes when checking next word (via trim)
	// Preserves original casing ("a"→"an", "A"→"An")
	wordsCopy := make([]string, len(words)) // Copy to avoid modifying caller slice
	copy(wordsCopy, words)
	vowels := "aeiouh" // Vowels + 'h' special case (an hour)

	for i := 0; i < len(wordsCopy); i++ {
		wordLower := strings.ToLower(wordsCopy[i])
		if wordLower != "a" && wordLower != "an" {
			continue // Not an article
		}

		// Find next non-empty word (skip punct/quotes)
		nextWordLower := ""
		for j := i + 1; j < len(wordsCopy); j++ {
			// Trim punct/quotes/punctuation from candidate next word
			trimmed := strings.Trim(wordsCopy[j], ".,!?:;'\\\"")
			if trimmed != "" {
				nextWordLower = strings.ToLower(trimmed)
				break
			}
		}
		if nextWordLower == "" {
			continue // No next word
		}

		nextStartsVowel := strings.ContainsRune(vowels, rune(nextWordLower[0]))

		// Apply article correction
		if nextStartsVowel {
			// Change to "an/An"
			if wordsCopy[i] == "a" {
				wordsCopy[i] = "an"
			} else if wordsCopy[i] == "A" {
				wordsCopy[i] = "An"
			}
		} else {
			// Change to "a/A"
			if wordsCopy[i] == "an" {
				wordsCopy[i] = "a"
			} else if wordsCopy[i] == "An" {
				wordsCopy[i] = "A"
			}
		}
	}
	return wordsCopy
}

// Milestone 8: Full pipeline // strings.Fields + transforms + join
func processText(text string) string {
	// 1. Split into lines
	lines := strings.Split(text, "\n")
	var outLines []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			outLines = append(outLines, "") // Preserve empty
			continue
		}

		// Text parsing/tokenize (custom for punct/parens/quotes)
		tokens := tokenize(line)

		// 2. Modifiers (Milestone 3/4)
		tokens = processModifiers(tokens)

		// Step 3: Fix articles "a" -> "an" before quotes so sees clean next words (e.g. "Apple")
		tokens = fixArticles(tokens)
		// Step 4: Handle single quotes: attach opening to next word, closing to previous
		tokens = fixQuotes(tokens)
		// Step 5: Attach all punctuation (!? . ,) to previous word (no space)
		tokens = fixPunctuation(tokens)

		// Step 6: Re-join tokens to single line with spaces
		outLines = append(outLines, strings.Join(tokens, " "))
	}

	return strings.Join(outLines, "\n")
}

// Custom tokenizer: splits on space but keeps:
// - Full parentheticals like (bin) or (up, 2) intact (including spaces)
// - Groups consecutive punctuation e.g. !? → "!?"
// - Standalone single quotes ' as separate tokens
func tokenize(s string) []string {
	var result []string
	var current string
	inParen := false  // Track if inside ( ... )
	punct := ".,!?:;" // Punctuation chars to group

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == '(' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			inParen = true
			current = "("
			continue
		}

		if c == ')' {
			current += ")"
			result = append(result, current)
			current = ""
			inParen = false
			continue
		}

		if inParen {
			current += string(c) // Keep all in paren incl spaces
			continue
		}

		if strings.ContainsRune(punct, rune(c)) {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			tmp := string(c)
			// Group consecutive punctuation e.g. !? → "!?"
			for i+1 < len(s) && strings.ContainsRune(punct, rune(s[i+1])) {
				tmp += string(s[i+1])
				i++ // Manual advance
			}
			result = append(result, tmp)
			continue
		}

		if c == '\'' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			result = append(result, "'") // Standalone quote token
			continue
		}

		if c == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			continue
		}

		current += string(c)
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}
