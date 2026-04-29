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
10. [isLetter() - Letter Helper](#10-isletter---letter-helper)
11. [Pipeline Order Rationale](#11-pipeline-order-rationale)
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
)
```
**Deep dive**:
| Package | Critical Functions | Why Essential |
|---------|-------------------|--------------|
| `fmt` | `Println()` | CLI feedback |
| `os` | `Args` `ReadFile()` `WriteFile()` | CLI args, I/O |
| `strconv` | `Atoi()` `ParseInt(base,64)` `FormatInt()` | Number conversions (hex/bin) |
| `strings` | `Split()` `Fields()` `Join()` `Trim()` `ToUpper/Lower/Title()` `ContainsRune()` `HasPrefix/Suffix()` | ALL text processing |

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
		return
	}
```
**Deep dive**:
- `os.Args`: `[]string` (0=progname, 1=input, 2=output)
- `len(os.Args)`: slice length
- `!= 3`: Exactly prog + 2 files
- `fmt.Println`: stdout, auto-newline
- `return`: Early exit (implicit status 0 if success)

**Edge**: `go run main.go` → usage; `go run main.go in.txt` → usage.

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
		return
	}
```
**Go idiom**: **Always** check `err != nil` immediately.
- `err.Error()` auto-stringified
- Early return = no further execution
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
		return
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

**Comment**: "Handles simple/numbered mods; skip mod tokens (original logic preserved, switch for future)."

**Purpose**: Parse/dispatch `(up)`, `(hex)`, `(up, 3)` modifiers.

```go
func processModifiers(words []string) []string {
	result := []string{}
```
`result`: Build processed words (modifiers consumed, not output).

```go
	for i := 0; i < len(words); i++ {
		word := words[i]
```
**Traditional index loop**: Need `i+1` access for numbered mods.

**Numbered mods block**:
```go
		if word == "(up," || word == "(low," || word == "(cap," || word == "(bin," && i+1 < len(words) {
```
**Deep**:
- Exact prefix match: `strings.Fields` splits `(up, 2)` → `["(up,", "2)"]`
- `&& i+1 < len(words)`: Bounds check
- **No parens in check**: Fields handles it

```go
			nStr := strings.Trim(words[i+1], ".,!?:;)")
```
**Mastery**:
- `Trim(set)` removes **from both ends**
- `"2),!"` → `"2"`
- Handles punct after number

```go
			if n, err := strconv.Atoi(nStr); err == nil {
```
- `Atoi`: String → `int` (decimal)
- `err == nil`: Parse success

```go
				for j := 1; j <= n; j++ {
					target := len(result) - j
					if target >= 0 {
```
**Backward apply**:
- `j=1`: Last word (`len(result)-1`)
- `j=2`: Second last etc.
- `target >= 0`: Safety (n > available words)

```go
						switch word {
						case "(up,": result[target] = strings.ToUpper(result[target])
						case "(low,": result[target] = strings.ToLower(result[target])
						case "(cap,": result[target] = strings.Title(result[target]) // Title > manual cap
						case "(bin,": 
							if v, err := strconv.ParseInt(result[target], 2, 64); err == nil {
								result[target] = strconv.FormatInt(v, 10)
							}
```
**Transform mastery**:
- `ToUpper/ToLower`: Full case change
- `Title`: First letter cap per word (title case), rest lower
- **ParseInt(str, base, bitSize)**: `base=2` binary → `int64`
- `FormatInt(int64, base=10)`: Back to decimal string
- **64 bit**: Handles huge numbers (e.g., 64-bit bin)

```go
			i++ // Skip number token
			continue
```
**Skip**: `i++` consumes `"2)"`, `continue` next iteration.

**Simple mods**:
```go
		suffix := ""
		clean := word
		for len(clean) > 0 && strings.ContainsRune(".,!?:;", rune(clean[len(clean)-1])) {
			suffix = string(clean[len(clean)-1]) + suffix
			clean = clean[:len(clean)-1]
		}
```
**Trailing punct mastery**:
- Loop from **end only** (preserves opening `(`)
- `rune(byte)`: Unicode safe
- `"word!"` → clean="word", suffix="!"
- **Prepends** to suffix (order preserved)

```go
		switch clean {
		case "(hex)":
			if len(result) > 0 {
				if v, err := strconv.ParseInt(result[len(result)-1], 16, 64); err == nil {
					result[len(result)-1] = strconv.FormatInt(v, 10)
				}
				result[len(result)-1] += suffix
			}
			continue
```
- **Last word**: Modifiers apply to previous
- `base=16`: Hex → decimal
- `+= suffix`: Reattach punct (`"FF!"` → `"255!"`)

```go
		case "(bin)", "(up)", "(low)", "(cap)":
			if len(result) > 0 {
				switch clean {
				case "(bin)": ParseInt(...,2,64) → FormatInt(10)
				case "(up)": ToUpper(last)
				case "(low)": ToLower(last)
				case "(cap)": Title(last)
				}
				result[len(result)-1] += suffix
			}
			continue
```
Same pattern, unified switch.

```go
		result = append(result, word)
```
**Normal word**: Pass through.

**Genius**: Modifiers consumed, transform prior words, punct preserved.

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
		if strings.HasPrefix(word, "'") && len(word) > 1 {
			quoteOpen = true
		}
		if strings.HasSuffix(word, "'") && len(word) > 1 {
			quoteOpen = false
		}
```
**Attached quotes**:
- `'word`: Prefix → open
- `word'`: Suffix → close
- `len > 1`: Ignore lone `'`

```go
		result = append(result, word)
```
Normal words pass through.

**State machine genius**: Handles all `' ` positions correctly.

---

## 7. fixPunctuation() - Punctuation Attacher

**Comment**: "Leading punct split/attach; isPunct for groups (original ContainsRune loop superior)."

**Purpose**: `hello ,world!` → `hello,world!`

```go
func fixPunctuation(words []string) []string {
	result := []string{}
	puncs := ".,!?:;"
```
**puncs**: Target chars (compact rune set).

**Leading punct**:
```go
	for _, word := range words {
		if len(word) > 1 && strings.ContainsRune(puncs, rune(word[0])) && !isPunctuation(word) {
```
**3 guards**:
1. `len > 1`: Not pure punct
2. `word[0]` punct? (`rune()` safe)
3. `!isPunctuation()`: Mixed content

**Example**: `",world"` (leading comma, content).

```go
			pEnd := 0
			for pEnd < len(word) && strings.ContainsRune(puncs, rune(word[pEnd])) {
				pEnd++
```
**Count leading punct**:
- `",,world"` → pEnd=2
- `ContainsRune(str, rune)`: O(n) rune check

```go
			prefix := word[:pEnd]  // ",,"
			if len(result) > 0 {
				result[len(result)-1] += prefix  // Attach NO SPACE
```
**Attach to prev**: `"hello" + ",,"` → `"hello,,"`

```go
			result = append(result, word[pEnd:])  // "world"
```
Rest as new word.

**Pure punct**:
```go
		} else if isPunctuation(word) {
			if len(result) > 0 {
				result[len(result)-1] += word  // "!!!" → attach
```
**`"!!!"`**: All punct → attach.

```go
		} else {
			result = append(result, word)
		}
```
Normal.

**Genius**: Handles leading groups (`,,`) and pure punct perfectly.

---

## 8. isPunctuation() - Punctuation Detector

**Comment**: "All chars punct? (range rune > byte loop)."

```go
func isPunctuation(s string) bool {
	if len(s) == 0 { return false }
	puncs := ".,!?:;"
	for _, r := range s {
		if !strings.ContainsRune(puncs, r) { return false }
	}
	return true
}
```
**Deep**:
- `range string`: **Runes** (not bytes) — UTF-8 safe
- `"é!"` → 2 runes ('é','!')
- `ContainsRune(str, rune)`: Searches rune set
- **Early false**: Non-punct found → false
- `"!!!"` → true

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
			for letIdx < len(next) && !isLetter(next[letIdx]) {
				letIdx++
```
**Skip non-letters**: `"'Apple"` → skip `'` → 'A'

```go
			if letIdx < len(next) && strings.ContainsRune(vowels, rune(next[letIdx])) {
				rep := "an"
				if len(clean) > 0 && clean[0] >= 'A' && clean[0] <= 'Z' {
					rep = "An"
```
**Replace**:
- Vowel → "an"/"An" (case match)
- `words[i][:prefixLen] + rep`: `"'a"` → `"'an"`

**In-place modify**: Efficient.

---

## 10. isLetter() - Letter Helper

**Comment**: "Helper for article next-word first letter."

```go
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
```
**Simple byte ranges**:
- ASCII letters only (performance)
- `next[letIdx]`: Safe byte access (post-Trim, ASCII assumed)
- **No rune**: `byte` faster for ASCII check

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
| **Switch** | `switch clean` | Fallthrough OFF, exhaustive possible |
| **String slicing** | `word[:pEnd]` | Zero-copy views |
| **strconv** | `ParseInt(base,64)` | Signed int64, base 2/10/16 |
| **Fields vs Split** | `strings.Fields` | Whitespace normalize, no empties |

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
