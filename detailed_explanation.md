# Detailed Explanation of go-reloaded Code

This document explains **every line** of the code in extreme detail. I'll break it down section by section, explaining:
- **What the code does**
- **Why it's written this way**
- **Go language concepts used**
- **Package details (os, strconv, strings)**
- **How to write it yourself step-by-step**
- **Common pitfalls & alternatives**
- **Performance/ best practices**

Goal: After reading, you can **fully understand, modify, and recreate** this code.

## 1. Package Declaration
```
package main
```
- **Purpose**: Every Go file belongs to a **package**. `main` is special - creates **executable** binary.
- **How Go works**: `go run main.go` → compiles to binary with entry point `main()`.
- **Why `main`**: Only packages named `main` produce executables. Libraries use other names (e.g. `fmt`).
- **To recreate**: Always `package main` for CLI tools.
- **Pitfall**: Wrong package name → `go run` fails with \"can't load package\".

## 2. Imports
```
import (
    \"os\" 
    \"strconv\"
    \"strings\"
)
```
- **Purpose**: Import **standard library packages** (no external deps - production ready).
- **os package details**:
  - **Purpose**: Operating System interface (files, args, processes).
  - **Key functions used**:
    - `os.Args []string`: Command line arguments. Index 0=program name, 1+=user args.
    - `os.ReadFile(filename) ([]byte, error)`: Reads entire file into memory. Returns bytes + error.
    - `os.WriteFile(filename, data []byte, perm os.FileMode) error`: Writes bytes to file with permissions.
  - **FileMode 0644**: Octal = owner read/write (6), group/world read-only (4).
  - **To learn more**: `go doc os` or https://pkg.go.dev/os

- **strconv package details**:
  - **Purpose**: String conversions (string ↔ numbers).
  - **Key functions**:
    - `strconv.ParseInt(s string, base int, bitSize int) (int64, error)`: Parses string to int.
      - `base 16` = hex, `base 2` = binary, `bitSize 64` = int64 range.
    - `strconv.Atoi(s string) (int, error)`: String to int (base 10).
    - `strconv.FormatInt(n int64, base int) string`: Number to string.
  - **Error handling**: Always check `err != nil`.
  - **To learn**: `go doc strconv`

- **strings package details**:
  - **Purpose**: String manipulation (immutable, efficient).
  - **Key functions used**:
    - `strings.ToLower/ToUpper(s string) string`: Case change.
    - `strings.Trim(s, chars string) string`: Remove chars from start/end.
    - `strings.Split(s, sep string) []string`: Split by separator.
    - `strings.ReplaceAll(old, new string) string`: Global replace.
    - `strings.Contains(token, \",\") bool`: Substring check.
    - `strings.Fields(s string) []string`: Split on whitespace (not used - custom tokenize instead).
  - **Rune**: `rune(word[0])` = Unicode codepoint (UTF-8 safe).
  - **To learn**: `go doc strings`

**How to choose imports**: Only import what you use. Go compiler errors on unused.

## 3. main() Function - Entry Point
```
func main() {
    // 1. Check that exactly 2 arguments...
```
- **Purpose**: Program entry point. Runs first.
- **Line-by-line**:

```
if len(os.Args) != 3 {
```
  - `os.Args` is `[]string` slice.
  - `len(os.Args)` = total args count.
  - Expect **exactly 3**: [0]=\"go-reloaded\", [1]=input file, [2]=output file.
  - **Example**: `go run main.go input.txt out.txt` → len=3.

```
println(\"Usage: go run . <input_file> <output_file>\")
return
```
  - `println` = simple print with newline (like fmt.Println but no format).
  - `return` exits main early (program ends).

```
inFilename := os.Args[1]
outFilename := os.Args[2]
```
  - `:=` = **short declaration** (declare + assign, type inferred).
  - `string` type from `os.Args[]string`.

```
content, err := readFile(inFilename)
if err != nil {
    println(\"Error reading file:\", err.Error())
    return
}
```
  - Calls `readFile`, unpacks tuple `(string, error)`.
  - **Idiomatic Go error handling**: `if err != nil { handle; return }`.
  - `err.Error()` → string message.

```
result := processText(content)
```
  - Main transformation call. `processText` signature `string → string`.

```
err = writeFile(outFilename, result)
if err != nil { ... }
```
  - Reuse `err` var (shadowing ok).
  - **Error propagation**: Functions return `error`, caller checks.

```
println(\"File processing completed successfully.\")
```
  - Success feedback.

**How to write main yourself**:
1. Check `len(os.Args)`.
2. Extract args.
3. Read → process → write.
4. Always handle errors.

## 4. readFile Function
```
func readFile(filename string) (string, error) {
```
- **Signature**: Input `string`, outputs `(string, error)` tuple.
- **Named return?** No - uses naked returns.

```
data, err := os.ReadFile(filename)
```
  - `os.ReadFile` reads **whole file** into `[]byte` (binary).
  - **Why bytes?** Files are binary, text is interpretation.
  - **When fails**: File not found, permissions, disk full.

```
if err != nil {
    return "", err
}
```
  - **Early return pattern**: Handle error, skip rest.
  - `""` empty string placeholder.

```
return strings.ReplaceAll(string(data), \"\\r\", \"\"), nil
```
  - `string(data)`: `[]byte` → string (UTF-8).
  - `strings.ReplaceAll`: Remove all `\r` (Windows line endings → Unix `\n`).
  - `nil` = no error.

**Master tip**: `os.ReadFile` loads **entire** file RAM - good for small files (<1GB).

**Alternative**: `bufio.Scanner` for large/streaming.

## 5. writeFile
```
func writeFile(filename string, content string) error {
    return os.WriteFile(filename, []byte(content), 0644)
}
```
- `[]byte(content)`: string → bytes.
- `0644` = file perms (visible `ls -l`).
- Single line - concise!

**Permissions explained**:
- 0 = no perms, 4=read, 6=read+write, 7=read+write+execute.
- `0644` = owner: rw (6), others: r (4).

## 6. hexToDecimal & binToDecimal - Number Conversion Helpers
```
func hexToDecimal(s string) string {
    n, err := strconv.ParseInt(s, 16, 64)
```
- **Input**: String like \"1E\".
- **base 16**: Hex (0-9A-F).
- **bitSize 64**: Returns `int64` (-9e18 to 9e18).
- **Error cases**: \"ZZZ\" not hex → err.

```
if err != nil { return s }
return strconv.FormatInt(n, 10)
```
- Fail safe: unchanged.
- `FormatInt(n, 10)`: int64 → decimal string.

**Same for binary** (base 2).

**Why helpers?** Reusable, testable, single responsibility.

**strconv mastery**:
- `ParseUint` for unsigned.
- `bitSize 0` = smallest type.
- Always check error!

## 7. processModifiers - Heart of Transformations
**Complex but powerful** - handles 7 modifier types.

```
var result []string
for i := 0; i < len(words); i++ {
```
- Build **new slice** (immutable strings).
- `i` index loop (manual for skipping).

**Modifier logic** (common pattern):
1. Check `token == \"(hex)\"`.
2. Apply to `result[len-1]` (last word).
3. `continue` skip adding token.

**Simple modifiers**: `(up)/(low)/(cap)` → `applyCase`.

**Numbered (up, 2)**:
```
if strings.HasPrefix(token, \"(\") && strings.Contains(token, \",\") {
```
- `HasPrefix`: Starts with \"(\" ?
- `Contains`: Has comma?

```
clean := strings.Trim(strings.ReplaceAll(token, \" \"), \"()\")
parts := strings.Split(clean, \",\")
```
- `ReplaceAll`: Remove spaces → \"(up,2)\".
- `Trim`: Remove () → \"up,2\".
- `Split`: [\"up\",\"2)\"].

```
n, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(parts[1]), \")\"))

```
- `TrimSpace`: Remove whitespace.
- `TrimSuffix`: Remove trailing \")\" → \"2\".
- `Atoi`: String → int.

```
start := len(result) - n
if start < 0 { start = 0 } // Guard
```
- Apply to **last N words** (or all if N > len).
- Negative index prevent panic.

```
for j := start; j < len(result); j++ {
    result[j] = applyCase(...)
```
- **In-place modify** slice elements.

```
result = append(result, token) // Normal word
```
- **Grow slice** dynamically.

**Master slice tip**: Slices are **views** on arrays. `append` may reallocate.

## 8. applyCase - Case Helper
```
switch mod {
case \"(up)\": strings.ToUpper(s)
case \"(low)\": strings.ToLower(s)
case \"(cap)\": TitleCase logic
```
- `switch` efficient (constant time).
- **Cap logic**: lower all → Upper first char + rest.

**Rune conversion**: `string(lower[0])` because `lower[0]` is byte.

## 9. fixPunctuation
```
punctuationChars := \".,!?:;\"
if strings.ContainsRune(punctuationChars, rune(word[0]))
```
- `ContainsRune`: Check if byte is punct (UTF-8 safe).
- `rune(word[0])`: First char as Unicode.

```
result[len(result)-1] += word // Direct concat
```
- **No space** between word+punct.

## 10. fixQuotes - State Machine
```
quoteOpen := false
```
- **Finite state**: Open/closed.

**Opening**:
```
if !quoteOpen && wordsCopy[i] == \"'\" {
    wordsCopy[i+1] = \"'\" + wordsCopy[i+1]
```
- Prefix next word. Modify copy.

**Closing**:
```
result[len(result)-1] += \"'\"
```
- Suffix **last added** word.

**Genius**: Modifies copy for opening, result for closing → handles multi-word quotes.

## 11. fixArticles - Lookahead Logic
```
vowels := \"aeiouh\"
```
- Hardcoded set.

**Lookahead**:
```
trimmed := strings.Trim(wordsCopy[j], \".,!?:;'\\\\\"\")
```
- **Trim attached** punct/quotes → check inner word.
- `ContainsRune(vowels, rune(next[0]))`.

**Bi-directional**: a→an OR an→a.

## 12. processText - Pipeline Orchestrator
**Line-by-line processing**.

```
lines := strings.Split(text, \"\\n\")
```
- Preserve structure.

**Custom order**:
1. tokenize (special tokens)
2. modifiers (numbers/case)
3. **articles** (needs clean tokens)
4. quotes (attaches to words)
5. punctuation (final attach)

**Why this order?** Articles see plain \"Apple\" before quotes make \"'Apple'\".

## 13. tokenize - Custom Parser
**Manual character loop** - most complex.

**State machine** for parens:
```
if inParen { current += c } // Preserve spaces in (up, 2)
```

**Punct grouping**:
```
tmp := string(c)
for next is punct { tmp += next; i++ }
```
- **Manual i++ skips** grouped chars.

**Quote special**:
```
result = append(result, \"'\")
```
- Separate token for fixQuotes to find.

**Final append** for trailing.

---

## How to Build This Code Yourself (Step-by-Step)

1. **Start minimal**:
```
package main
import \"os\"
func main() { println(len(os.Args)) }
```

2. **Add file I/O**:
```
data, _ := os.ReadFile(os.Args[1])
os.WriteFile(os.Args[2], data, 0644)
```

3. **Add tokenize** (character by character).

4. **Add one modifier at a time**.

5. **Test each step**: `go test`, manual inputs.

## Go Best Practices Used
- **Error handling everywhere**.
- **Short declarations `:=`**.
- **Early returns**.
- **Slice building with append**.
- **Copy slices to avoid mutation**.
- **Constants for magic strings**.

## Common Mistakes Avoided
- **No `strings.Fields`**: Misses punct/parens.
- **Index bounds**: Guards everywhere.
- **UTF-8 safe**: Runes for chars.
- **Memory**: Whole file ok for small.

## Performance
- O(n) everything.
- No allocations in hot loops (reuse slices).

**Congratulations** - production-quality code!

Mastered packages? Run `go doc os strconv strings` daily.

