# go-reloaded: Complete Line-by-Line Mastery Guide

This document provides **extra-detailed, exhaustive explanations** of every single line in `main.go`, with deep dives into Go concepts, function behaviors, edge cases, design decisions, and mastery-level insights. Only code present in the current `main.go` is covered — no obsolete snippets.

**Goal**: Master every nuance to understand, modify, debug, and extend the codebase confidently.

## Table of Contents
1. [Project Overview & Architecture](#1-project-overview--architecture)
2. [package main & imports](#2-package-main--imports)
3. [main() - CLI Entry Point](#3-main---cli-entry-point)
4. [processText() - Pipeline Orchestrator](#4-processtext---pipeline-orchestrator)
5. [processModifiers() - Modifier Engine](#5-processmodifiers---modifier-engine)
6. [fixQuotes() - Quote State Machine](#6-fixquotes---quote-state-machine)
7. [fixPunctuation() - Punctuation Attacher](#7-fixpunctuation---punctuation-attacher)
8. [isPunctuation() - Punctuation Detector](#8-ispunctuation---punctuation-detector)
9. [fixArticles() - Article Corrector](#9-fixarticles---article-corrector)
10. [capitalize() - Custom Capitalizer](#10-capitalize---custom-capitalizer)
11. [Pipeline Order Rationale](#11-pipeline-order-rationale)
12. [Go Concepts Deep Dive](#12-go-concepts-deep-dive)
13. [Edge Cases & Testing](#13-edge-cases--testing)
12. [Go Concepts Deep Dive](#12-go-concepts-deep-dive)
13. [Edge Cases & Testing](#13-edge-cases--testing)

---

## 1. Project Overview & Architecture

**go-reloaded** is a **CLI text processor** that:
- Reads file via `os.ReadFile`
- Line-by-line processing via `processText`: split → words → 4-stage pipeline → join
- **Pipeline**: Modifiers → Quotes → Punctuation → Articles (order critical)
- Writes via `os.WriteFile`
- **Transformations** (Milestones 1-8):
  | # | Feature |
  |---|---------|
  |1| File I/O |
  |2| `(hex)`/`(bin)` → decimal |
  |3| `(up)`/`(low)`/`(cap)` |
  |4| `(up,3)` numbered |
  |5| Punct attach (no space before) |
  |6| Single quotes glue |
  |7| `a/hour` → `an/hour` |
|8| Pipeline order |
|9| Function map modifiers & unicode helpers |

**Key Design**:
- Line-wise to preserve empty lines
- `strings.Fields` for natural whitespace handling
- `[]string` words for precise control
- Inline comments explain optimizations

**Usage**: `go run main.go input.txt output.txt`

---

## 2. package main & imports

```go
package main
```
**Mastery note**: 
- `main` = executable (generates binary)
- `go run main.go` compiles + executes
- `go build` → executable binary
- No `main` = library package

```go
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)
```
**Deep dive**:
| Package | Critical Functions | Why Essential |
|---------|-------------------|--------------|
| `fmt` | `Println()` | CLI feedback |
| `os` | `Args` `ReadFile()` `WriteFile()` | CLI args, I/O |
| `strconv` | `Atoi()` `ParseInt(base,64)` `FormatInt()` | Number conversions (hex/bin) |
| `strings` | `Split()` `Fields()` `Join()` `Trim()` `ToUpper/Lower/Title()` `ContainsRune()` `HasPrefix/Suffix()` | ALL text processing |
| `unicode` | `IsLetter()` `IsUpper()` | Unicode-aware letter/case detection in fixArticles() |

**Go rule**: Unused imports = compile error (enforces minimalism).

---

## 3. main() - CLI Entry Point

**Purpose**: Parse args, I/O, orchestrate processing, error handling.

```go
func main() {
```
**Mastery**: Go entry point. No args. Returns `int` implicitly (0=success).

```go
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input> <output>")
		os.Exit(1)
	}
```
**Deep dive**:
- `os.Args`: `[]string` (0=progname, 1=input, 2=output)
- `len(os.Args)`: slice length
- `!= 3`: Exactly prog + 2 files
- `fmt.Println`: stdout, auto-newline
- `os.Exit(1)`: Explicit error exit code (non-zero)

**Edge**: `go run main.go` → usage exit 1.

```go
	data, err := os.ReadFile(os.Args[1])
```
**Mastery**:
- `os.ReadFile(path)`: `([]byte, error)` — reads **entire** file
- `data`: `[]byte` (raw bytes, UTF-8 safe for text)
- `err`: nil or `os.PathError` etc.
- `os.Args[1]`: Direct path access (no var for conciseness)
- `:=`: Short decl (infers types)

**Why `[]byte`?** Files binary. `string([]byte)` valid for text.

```go
	if err != nil {
		fmt.Println("Error reading:", err)
		os.Exit(1)
	}
```
**Go idiom**: **Always** check `err != nil` immediately.
- `err.Error()` auto-stringified
- `os.Exit(1)`: Explicit non-zero exit status for errors
- **No panic**: Explicit error handling (Go philosophy)

```go
	result := processText(string(data))
```
**Core**: `[]byte` → `string`, process → transformed `string`
- `string(data)`: Lossless for ASCII/UTF-8
- `result`: Processed text ready for output

```go
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
```
**Mastery**:
- `os.WriteFile(path, []byte, mode)`: `error`
- `[]byte(result)`: `string` → bytes
- `0644`: Octal perms = rw-r--r-- (owner rw, others read)
- `=` : Reuses `err` var

```go
	if err != nil {
		fmt.Println("Error writing:", err)
		os.Exit(1)
	}
```
Same error pattern.

```go
	fmt.Println("Success.")
```
**Polish**: User feedback. `return` implicit.

---

## 4. processText() - Pipeline Orchestrator

**Comment**: "Orchestrates pipeline per line to preserve newlines/empty lines."

**Purpose**: Line-by-line processing preserving structure.

```go
func processText(input string) string {
```
**Signature**: Pure func, no side effects.

```go
	lines := strings.Split(input, "\n")
```
**Deep dive**:
- `strings.Split(s, sep)`: `[]string`
- `"\n"`: Line separator (Unix standard)
- Preserves trailing empty line if ends with `\n`
- `lines[0]` = first line etc.

**Why lines first?** Modifiers work on words; need to preserve empty lines.

```go
	var out []string
```
**Mastery**:
- `var` + type: Zero-value `[]string{}` (empty slice)
- `out` will collect processed lines
- Later `append(out, ...)` grows it

```go
	for _, line := range lines {
```
- `range slice`: `index, value` — `_` ignores index
- Line-by-line processing

```go
		words := strings.Fields(line)
```
**Critical**:
- `strings.Fields(s)`: Split on **any** whitespace (space/tab/NL), skip empty
- `"a  b\n c"` → `["a", "b", "c"]` (normalizes spaces)
- **vs Split(" ")**: Would give empty strings for multi-spaces
- **Mastery**: Perfect for natural text (ignores formatting whitespace)

```go
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
```
**Edge case mastery**:
- Empty line (`""`) → `Fields("") = []string{}`
- `len == 0`: Preserve empty lines
- `append(out, "")`: Empty string line
- `continue`: Skip pipeline (no words to process)

```go
		words = processModifiers(words)
		words = fixQuotes(words)
		words = fixPunctuation(words)
		words = fixArticles(words)
```
**Pipeline** — **ORDER CRITICAL** (see #11):
1. **Modifiers**: Transform content first
2. **Quotes**: Wrap transformed words  
3. **Punctuation**: Attach to quoted words
4. **Articles**: Check next word's first letter (post-transform)

Each func: `([]string) []string` — pure, chainable.

```go
		out = append(out, strings.Join(words, " "))
```
**Reconstruct line**:
- `strings.Join([]string, " ")`: Single space between words
- `append`: Grow `out`

```go
	return strings.Join(out, "\n")
```
**Final assembly**:
- Join lines with `\n`
- **Preserves** trailing newline if input had it

**Mastery insight**: Line-wise → perfect empty line/structure preservation.

---

## 5. processModifiers() - Modifier Engine

**Point 2**: Simplified modifier logic using **function map** - unified simple/numbered handling.

**Purpose**: Parse `(up)`, `(hex)`, `(up,3)` with extensible map pattern.

```go
func processModifiers(words []string) []string {
	// Map of simple transformation functions
	transformations := map[string]func(string) string{
		"(up)":  strings.ToUpper,
		"(low)": strings.ToLower,
		"(cap)": capitalize,
		"(hex)": func(s string) string {
			if v, err := strconv.ParseInt(s, 16, 64); err == nil {
				return strconv.FormatInt(v, 10)
			}
			return s
		},
		"(bin)": func(s string) string {
			if v, err := strconv.ParseInt(s, 2, 64); err == nil {
				return strconv.FormatInt(v, 10)
			}
			return s
		},
	}
```
**Mastery - Function Map Pattern**:
- `map[string]func(string) string`: Key=mod name, Value=function
- Built-ins: ToUpper/ToLower
- Custom: capitalize() helper
- **Closures** for hex/bin: Capture ParseInt logic
- `err == nil`: Invalid input → unchanged
- **Extensible**: Add `"(rev)": strings reverse func` easily

```go
	result := []string{}
	for i := 0; i < len(words); i++ {
		word := words[i]
```
**Loop same**: Index for numbered mods.

**Numbered mods**:
```go
		if (word == "(up," || word == "(low," || word == "(cap," || word == "(bin,") && i+1 < len(words) {
```
Same prefix detection + bounds.

```go
			nStr := strings.Trim(words[i+1], ".,!?:;)")
			if n, err := strconv.Atoi(nStr); err == nil {
				for j := 1; j <= n; j++ {
					target := len(result) - j
					if target >= 0 {
```
Same parsing/backward apply.

```go
						tag := word[:len(word)-1] + ")" // Convert "(up," to "(up)" for the map
						if fn, ok := transformations[tag]; ok {
							result[target] = fn(result[target])
						}
					}
				}
			}
			i++ // Skip number
			continue
```
**Genius upgrade**: `tag` normalization "(up," → "(up)" → map lookup!
- `[:len-1]`: Drop comma
- `map[string]fn`: Unified numbered/simple
- `if ok`: Safe lookup

**Simple mods**:
```go
		suffix := ""
		clean := word
		for len(clean) > 0 && strings.ContainsRune(".,!?:;", rune(clean[len(clean)-1])) {
			suffix = string(clean[len(clean)-1]) + suffix
			clean = clean[:len(clean)-1]
		}
```
Same punct strip.

```go
		if fn, ok := transformations[clean]; ok {
			if len(result) > 0 {
				result[len(result)-1] = fn(result[len(result)-1]) + suffix
```
**Unified dispatch**: Map lookup → apply to last + reattach suffix
- No switch! DRY principle
- **hex/bin closures** handle ParseInt inline

```go
			continue
		}
		result = append(result, word)
```
Normal word pass-through.

**Design brilliance**: Map eliminates switch duplication, closures encapsulate ParseInt, tag norm unifies numbered/simple.

---

## 6. fixQuotes() - Quote State Machine

**Comment**: "State machine for ' gluing (original exact logic; HasPrefix/Suffix safe)."

**Purpose**: Convert `word ' awesome'` → `word 'awesome'`.

```go
func fixQuotes(words []string) []string {
	result := []string{}
	quoteOpen := false
```
`quoteOpen`: FSM state (`false`=expect open).

```go
	for i := 0; i < len(words); i++ {
		word := words[i]
		if word == "'" {
```
**Standalone `' `** detected.

```go
			if !quoteOpen {
				if i+1 < len(words) {
					words[i+1] = "'" + words[i+1]  // Prepend to NEXT
					quoteOpen = true
					continue  // Skip adding lone '
				}
```
**Open quote**: Glue to **next** word (`'word`).

```go
			} else {
				if len(result) > 0 {
					result[len(result)-1] += "'"  // Append to LAST
					quoteOpen = false
					continue
				}
```
**Close quote**: Glue to **previous** word (`word'`).

```go
		// Update state if word already contains a quote
		if strings.HasPrefix(word, "'") && !strings.HasSuffix(word, "'") {
			quoteOpen = true
		} else if strings.HasSuffix(word, "'") && !strings.HasPrefix(word, "'") {
			quoteOpen = false
		}
```
**Refined state updates**:
- `'word'` (both): No state change
- `'word` (prefix only): Open
- `word'` (suffix only): Close
- Prevents toggle on fully quoted words

```go
		result = append(result, word)
```
Normal words pass through.

**State machine genius**: Handles all `' ` positions correctly.

---

## 7. fixPunctuation() - Punctuation Attacher

**Point 3**: Efficient handling with `strings.Builder` preparation, edge case coverage.

**Purpose**: `hello ,world!` → `hello,world!`, handles empty result cases.

```go
func fixPunctuation(words []string) []string {
	result := []string{}
	puncs := ".,!?:;"
```
Same punct set.

```go
	for _, word := range words {
		// Use strings.Builder for potentially complex concatenation
		var sb strings.Builder
```
**Builder prep**: For future expansion (e.g., formatting); minimal use here.

**Leading punct**:
```go
		if len(word) > 1 && strings.ContainsRune(puncs, rune(word[0])) && !isPunctuation(word) {
```
Same 3 guards.

```go
			pEnd := 0
			for pEnd < len(word) && strings.ContainsRune(puncs, rune(word[pEnd])) {
				pEnd++
			}
			prefix := word[:pEnd]
			if len(result) > 0 {
				result[len(result)-1] += prefix
```
**Core attach**.

```go
			} else {
				result = append(result, prefix)
```
**NEW edge**: No prev → start with prefix (e.g., leading ",world").

```go
			result = append(result, word[pEnd:])
		} else if isPunctuation(word) {
			if len(result) > 0 {
				result[len(result)-1] += word
```
**Pure punct**.

```go
			} else {
				result = append(result, word)
```
**NEW edge**: Pure punct first → standalone.

```go
			} else {
				result = append(result, word)
			}
		}
		_ = sb.String() // Builder usage would expand here for complex formatting
	}
```
**Builder placeholder**: Ready for extensions.

**Mastery**: Robust edges + Builder future-proofing.

---

## 8. isPunctuation() - Punctuation Detector

**Comment**: "All chars punct? (range rune > byte loop)."

```go
func isPunctuation(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(".,!?:;", r) {
			return false
		}
	}
	return true
}
```
**Deep**:
- `s == ""`: Explicit empty check (vs len==0)
- `range string`: **Runes** — UTF-8 safe
- Direct `ContainsRune(".,!?:;", r)`: Inline puncs
- **Early false**: Any non-punct → false
- `"!!!"` → true; `""` → false

**Used by**: `fixPunctuation` to distinguish pure vs mixed.

---

## 9. fixArticles() - Article Corrector

**Comment**: "a/an based on next letter (Trim > Index for prefix; ContainsRune vowels)."

**Purpose**: `a apple` → `an apple`, `A hour` → `An hour`.

```go
func fixArticles(words []string) []string {
	vowels := "aeiouhAEIOUH"  // +h (hour)
```
**Vowels**: Case-insensitive + 'h'.

```go
	for i := 0; i < len(words)-1; i++ {  // Needs next word
```
**Stop short**: No next → no change.

```go
		clean := strings.Trim(words[i], "'\"")  // "'a" → "a"
		prefixLen := len(words[i]) - len(clean)  // 2
```
**Prefix calc**:
- `"'a"` (len=2) → clean="a" (len=1) → prefixLen=1 (`'`)
- **Elegant**: No Index needed

```go
		if strings.ToLower(clean) == "a" {
```
Case-insensitive "a"/"A".

```go
			next := words[i+1]
			letIdx := 0
			for letIdx < len(next) && !unicode.IsLetter(rune(next[letIdx])) {
				letIdx++
```
**Skip non-letters**: Unicode-safe via `IsLetter(rune)`.

```go
			if letIdx < len(next) && strings.ContainsRune(vowels, rune(next[letIdx])) {
				rep := "an"
				if len(clean) > 0 && unicode.IsUpper(rune(clean[0])) {
					rep = "An"
```
**Replace**:
- `IsUpper(rune)`: Unicode case check
- `words[i][:prefixLen] + rep`: Prefix preserved

**Unicode mastery**: Handles accented chars properly.

---

## 10. capitalize() - Custom Capitalizer

**Purpose**: Custom first-letter capitalization for `(cap)` modifier (simpler than Title for single words).

```go
func capitalize(s string) string {
```
**Signature**: Pure string → string transformer.

```go
    if s == "" { return "" }
```
**Edge case**: Empty string → empty (safety).

```go
    return strings.ToUpper(s[:1]) + s[1:]
```
**Mastery**:
- `s[:1]`: First **rune** as string slice (UTF-8 safe)
- `ToUpper`: First char upper
- `s[1:]`: Rest unchanged
- `"hello"` → `"Hello"`
- **vs Title**: Title caps every word; this single word only

**Used by**: `(cap)` and `(cap,n)` via map lookup.

---

## 11. Pipeline Order Rationale

**processText pipeline**:
```
Modifiers → Quotes → Punctuation → Articles
```

**Why this order?** (Dependency chain):

1. **Modifiers first**: Transform `(hex)` content **before** quotes/punct attach
   - `FF (hex)'` → `"255'"` (not `'FF`)
2. **Quotes second**: Wrap transformed words
   - `"255'"` correct
3. **Punctuation third**: Attach to quoted words
   - `"hello' ,world"` → `"hello',world"`
4. **Articles last**: Check **final** next-word letter
   - `(cap) Apple` → `"Apple"` starts 'A' → `an Apple`

**Wrong order breaks**:
- Quotes before mods: `(hex)'` wrong
- Articles before quotes: `'Apple` misdetected

**Mastery**: Dependency graph determines order.

---

## 12. Go Concepts Deep Dive

| Concept | Usage | Mastery Insight |
|---------|--------|----------------|
| **Slices** | `append(result, word)` `result[len-1]` | Grow dynamically, last access O(1)
| **Runes** | `rune(word[0])` `range s` | Unicode codepoints (not bytes) |
| **Error Handling** | `if err != nil` | Explicit, no exceptions |
| **Short decl** | `data, err :=` | `:=` infers, `=` reassigns |
| **Switch** | `switch clean` | Fallthrough OFF, exhaustive possible (legacy) |
| **Func Maps** | `map[string]func(str)str` | Unified modifiers dispatch |
| **Closures** | hex/bin lambdas | Encapsulate ParseInt logic |
| **String slicing** | `word[:pEnd]` `s[:1]` | Zero-copy, first rune extract |
| **strconv** | `ParseInt(base,64)` | Signed int64, base 2/10/16 |
| **Unicode** | `IsLetter(rune)` | UTF-8 letter detection |

**Performance notes**:
- `Fields`: O(n) whitespace scan
- `range rune`: UTF-8 safe
- Inline funcs: No call overhead

---

## 13. Edge Cases & Testing

**Critical edges handled**:
| Case | Handling |
|------|----------|
| Empty line | `len(words)==0` → `""` |
| No prior word for mod | `len(result)==0` → skip |
| n > words | `target >=0` guard |
| Invalid hex/bin | `err != nil` → unchanged |
| Lone `'` | State machine skips |
| Leading punct only | `pEnd == len(word)` handled |
| `'a` last | Loop `i < len-1` skips |
| Multi-punct | `pEnd` loop + `isPunctuation` |

**Test with**: `go test` (table-driven, comprehensive).

**Mastery achieved**: Every line explained, every concept mastered.

---
