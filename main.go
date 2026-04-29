package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input> <output>")
		return
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	result := processText(string(data))
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
	fmt.Println("Success.")
}

// processText: Orchestrates pipeline per line to preserve newlines/empty lines.
func processText(input string) string {
	lines := strings.Split(input, "\n")
	var out []string
	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		words = processModifiers(words)
		words = fixQuotes(words)
		words = fixPunctuation(words)
		words = fixArticles(words)
		out = append(out, strings.Join(words, " "))
	}
	return strings.Join(out, "\n")
}

// processModifiers: Handles simple/numbered mods; skip mod tokens (original logic preserved, switch for future).
func processModifiers(words []string) []string {
	result := []string{}
	for i := 0; i < len(words); i++ {
		word := words[i]
		// Numbered: exact prefix match (Fields splits (up, 2) -> "(up,", "2)")
		if word == "(up," || word == "(low," || word == "(cap," || word == "(bin," && i+1 < len(words) {
			nStr := strings.Trim(words[i+1], ".,!?:;)")
			if n, err := strconv.Atoi(nStr); err == nil {
				for j := 1; j <= n; j++ {
					target := len(result) - j
					if target >= 0 {
						switch word {
						case "(up,":
							result[target] = strings.ToUpper(result[target])
						case "(low,":
							result[target] = strings.ToLower(result[target])
						case "(cap,":
							result[target] = strings.Title(result[target]) // Title > manual cap (handles multi-word).
						case "(bin,":
							if v, err := strconv.ParseInt(result[target], 2, 64); err == nil {
								result[target] = strconv.FormatInt(v, 10)
							}
						}
					}
				}
			}
			i++
			continue
		}
		// Simple mods: trailing punct strip (ContainsRune > Trim ends), apply/skip.
		suffix := ""
		clean := word
		for len(clean) > 0 && strings.ContainsRune(".,!?:;", rune(clean[len(clean)-1])) {
			suffix = string(clean[len(clean)-1]) + suffix
			clean = clean[:len(clean)-1]
		}
		switch clean {
		case "(hex)":
			if len(result) > 0 {
				if v, err := strconv.ParseInt(result[len(result)-1], 16, 64); err == nil {
					result[len(result)-1] = strconv.FormatInt(v, 10)
				}
				result[len(result)-1] += suffix
			}
			continue
		case "(bin)", "(up)", "(low)", "(cap)":
			if len(result) > 0 {
				switch clean {
				case "(bin)":
					if v, err := strconv.ParseInt(result[len(result)-1], 2, 64); err == nil {
						result[len(result)-1] = strconv.FormatInt(v, 10)
					}
				case "(up)":
					result[len(result)-1] = strings.ToUpper(result[len(result)-1])
				case "(low)":
					result[len(result)-1] = strings.ToLower(result[len(result)-1])
				case "(cap)":
					result[len(result)-1] = strings.Title(result[len(result)-1])
				}
				result[len(result)-1] += suffix
			}
			continue
		}
		result = append(result, word)
	}
	return result
}

// fixQuotes: State machine for ' gluing (original exact logic; HasPrefix/Suffix safe).
func fixQuotes(words []string) []string {
	result := []string{}
	quoteOpen := false
	for i := 0; i < len(words); i++ {
		word := words[i]
		if word == "'" {
			if !quoteOpen {
				if i+1 < len(words) {
					words[i+1] = "'" + words[i+1]
					quoteOpen = true
					continue
				}
			} else {
				if len(result) > 0 {
					result[len(result)-1] += "'"
					quoteOpen = false
					continue
				}
			}
		}
		if strings.HasPrefix(word, "'") && len(word) > 1 {
			quoteOpen = true
		}
		if strings.HasSuffix(word, "'") && len(word) > 1 {
			quoteOpen = false
		}
		result = append(result, word)
	}
	return result
}

// fixPunctuation: Leading punct split/attach; isPunct for groups (original ContainsRune loop superior).
func fixPunctuation(words []string) []string {
	result := []string{}
	puncs := ".,!?:;"
	for _, word := range words {
		if len(word) > 1 && strings.ContainsRune(puncs, rune(word[0])) && !isPunctuation(word) {
			pEnd := 0
			for pEnd < len(word) && strings.ContainsRune(puncs, rune(word[pEnd])) {
				pEnd++
			}
			prefix := word[:pEnd]
			if len(result) > 0 {
				result[len(result)-1] += prefix
			} else {
				result = append(result, prefix)
			}
			result = append(result, word[pEnd:])
		} else if isPunctuation(word) {
			if len(result) > 0 {
				result[len(result)-1] += word
			} else {
				result = append(result, word)
			}
		} else {
			result = append(result, word)
		}
	}
	return result
}

// isPunctuation: All chars punct? (range rune > byte loop).
func isPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	puncs := ".,!?:;"
	for _, r := range s {
		if !strings.ContainsRune(puncs, r) {
			return false
		}
	}
	return true
}

// fixArticles: a/an based on next letter (Trim > Index for prefix; ContainsRune vowels).
func fixArticles(words []string) []string {
	vowels := "aeiouhAEIOUH"
	for i := 0; i < len(words)-1; i++ {
		clean := strings.Trim(words[i], "'\"")
		prefixLen := len(words[i]) - len(clean)
		if strings.ToLower(clean) == "a" {
			next := words[i+1]
			letIdx := 0
			for letIdx < len(next) && !isLetter(next[letIdx]) {
				letIdx++
			}
			if letIdx < len(next) && strings.ContainsRune(vowels, rune(next[letIdx])) {
				rep := "an"
				if len(clean) > 0 && clean[0] >= 'A' && clean[0] <= 'Z' {
					rep = "An"
				}
				words[i] = words[i][:prefixLen] + rep
			}
		}
	}
	return words
}

// isLetter: Helper for article next-word first letter.
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
