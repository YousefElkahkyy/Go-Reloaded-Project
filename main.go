package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*

Milestone 1 — Read Input, Write Output

    How many command-line arguments does your program expect? What happens if the user provides the wrong number?  The program expects 2 command-line arguments: the input file and the output file. If the user provides the wrong number of arguments, the program will print a usage message and exit without processing any files.
    What type does os.ReadFile return, and what do you need to do before you can work with it as text? os.ReadFile returns a byte slice ([]byte) and an error. To work with it as text, you need to convert the byte slice to a string using string(date).
    What does the third argument of os.WriteFile control? The third argument of os.WriteFile controls the file permissions for the newly created file. In this case, 0644 means that the owner can read and write the file, while others can only read it.

Milestone 2 — Number Conversions: (hex) and (bin)

    What function converts a string like "42" from base 16 to a decimal integer? You can use strconv.ParseInt with a base of 16 to convert a hexadecimal string to a decimal integer. For example: strconv.ParseInt("42", 16, 64). And for binary, you can use strconv.ParseInt with a base of 2. For example: strconv.ParseInt("1010", 2, 64).
    What base value do you pass for hexadecimal? For binary? For hexadecimal, you pass a base value of 16. For binary, you pass a base value of 2.
    What should your function return if the input is not a valid number? If the input is not a valid number, the function should return the original string unchanged.

Milestone 3 — Case Modifiers: (up), (low), (cap)

    After splitting the text into words, how do you detect that a word is (up), (low), or (cap)? You can check if the current word matches any of the modifiers (up), (low), or (cap) using a simple if statement or a switch case.
    When you encounter a modifier, which word do you transform — and where is it relative to the modifier in your result slice? When you encounter a modifier, you should transform the last word that was added to the result slice, which is the word immediately before the modifier in the original text.
    What does capitalize mean exactly? What should happen to the letters after the first one? Capitalizing a word means converting the first letter to uppercase and the rest of the letters to lowercase. For example, "hELLO" would become "Hello" when capitalized.

Milestone 4 — Numbered Modifiers: (up, 2), (cap, 4)

    In your word slice, how does (up, appear as a token? What comes right after it? In the word slice, (up, would be a token that appears as "(up," and it would be followed by a number token like "2)" which indicates how many words to modify.
    How do you extract the number from a token like "2)"? You can use strings.Trim to remove the closing parenthesis and any punctuation, and then use strconv.Atoi to convert the remaining string to an integer. For example: nStr := strings.Trim("2)", ".,!?:;)") followed by n, err := strconv.Atoi(nStr).
    What should happen if the number is larger than the number of words already in result?  If the number is larger than the number of words already in result, you should only modify as many words as are available. For example, if you have (up, 5) but there are only 3 words before it, you would only modify those 3 words.

Milestone 5 — Punctuation Spacing

    How do you detect that a word is a punctuation token? You can check if the current word is one of the punctuation marks (.,!?:;) using a simple if statement. You might also want to check if the entire word consists of punctuation characters to handle cases where punctuation is attached to modifiers.
    Instead of adding a punctuation token as a new element, what do you do with it? Instead of adding a punctuation token as a new element in the result slice, you should attach it to the last word that was added to the result slice. This means concatenating the punctuation token to the end of the last word without adding a space.
    What does "attach to the previous word" look like in terms of your result slice?  If your result slice currently has ["Hello", "world"] and you encounter a punctuation token "!", instead of adding it as a new element, you would modify the last element to be "world!" resulting in ["Hello", "world!"].

Milestone 6 — Single Quote Formatting

    How do you find the opening ' and the matching closing ' in your word slice? to find the opening and closing single quotes, you can maintain a boolean variable (e.g., quoteOpen) that toggles when you encounter a standalone single quote token. When you see a single quote, if quoteOpen is false, it means it's an opening quote, and you set quoteOpen to true. If quoteOpen is true, it means it's a closing quote, and you set quoteOpen back to false.
    What does "attach the opening quote to the first word inside" look like in your result slice? If you encounter an opening quote and the next word is "hello", you would modify the next word to be "'hello" and not add the standalone quote to the result slice. So if your result slice was ["Hello"] and you encounter a standalone quote followed by "world", you would modify "world" to be "'world" and your result slice would become ["Hello", "'world"].
    How do you handle multiple quote pairs in the same line? You would continue to toggle the quoteOpen variable each time you encounter a standalone single quote token. This way, you can correctly identify multiple pairs of quotes in the same line and attach them to the appropriate words.

Milestone 7 — Article Correction: a → an

    How do you check the first character of the next word? You can check the first character of the next word by accessing the first index of the string (e.g., next[0]) and checking if it is a vowel (a, e, i, o, u, h) using strings.ContainsRune.
    How do you preserve the original casing — a → an, A → An? You can check the original casing of the article "a" by comparing it to "a" and "A". If the original article is lowercase "a", you would change it to "an". If the original article is uppercase "A", you would change it to "An". This way, you preserve the original casing while correcting the article.
    What happens if a is the last word in the text with no word after it? If "a" is the last word in the text with no word after it, you would not change it to "an" since there is no next word to check for a vowel. In this case, you would simply leave it as "a".

Milestone 8 — Connect the Pipeline

    What is the correct order to run the transformations? Does order matter? Yes the order of the transformations matters because some transformations depend on the results of previous ones. For example, you should fix Quotes before fixing Punctuation Spacing because if you fix the punctuation spacing first, you might end up with cases where a quote is attached to a punctuation mark, which could make it harder to correctly identify and format the quotes. By fixing the punctuation spacing first, you ensure that any punctuation is properly attached to the words before you handle the quotes, allowing for more accurate formatting of both punctuation and quotes in the final output.
    Why use strings.Fields to split instead of strings.Split(text, " ")? Using strings.Fields to split the text is generally better for processing natural language because it automatically handles multiple consecutive whitespace characters and treats them as a single separator. This means that if there are extra spaces between words, they won't result in empty strings in the resulting slice. On the other hand, using strings.Split(text, " ") would create empty string elements in the slice for each extra space, which can complicate processing and may require additional handling to filter out those empty strings. However, if you need to preserve the exact spacing (including multiple spaces), you might choose to use strings.Split instead and handle the empty strings accordingly.


So the arrange of the transformations in the pipeline is as follows:
1. Process Modifiers (Milestones 2, 3, 4)
2. Fix Quotes (Milestone 6)
3. Fix Punctuation Spacing (Milestone 5)
4. Fix Articles (Milestone 7)

words = processModifiers(words)
words = fixQuotes(words)
words = fixPunctuationSpacing(words)
words = fixArticles(words)

*/

func main() {
	arguments := os.Args

	// 1. Check that exactly 2 arguments were provided (input and output file)
	if len(arguments) != 3 {
		// Print usage message
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}

	inPutFile := arguments[1]
	outPutFile := arguments[2]

	// 2. Read the input file into a string
	content, err := readFile(inPutFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	//     // 3. (Transformations will go here in later milestones)

	//    /* // This will delete the whitespace, but it will also delete the newlines, which we want to preserve. We will need to handle this more carefully in later milestones.
	//     // split The raw content into a slice of words
	//     words := strings.Fields(content)
	//     // Process the modifiers in the words slice
	//     processedWords := processModifiers(words)
	//     // Join the processed words back into a single string
	//     content = strings.Join(processedWords, " ")
	//     */

	//     // better way to split the content into words while preserving newlines is to use strings.Split(content, " ") instead of strings.Fields(content). This will give us a slice of words that includes empty strings for consecutive spaces and newlines, allowing us to preserve the original formatting more accurately. We will need to handle the modifiers and transformations carefully to ensure we don't lose the newlines in the process.
	//     lines := strings.Split(content, "\n")
	//     var finalLines []string // This will hold the processed lines
	//     for _,line := range lines {
	//         // split the line into words
	//         /* // if we use Split in side the loop, t treats multiple spaces as "empty strings" in your slice. If you have "word    (up)", your slice looks like ["word", "", "", "", "(up)"], which can be useful for preserving the exact spacing. However, it also means you have to handle those empty strings in your processing logic, which can be a bit more complex. On the other hand, using Fields would give you ["word", "(up)"], which is simpler to process but loses the original spacing. Depending on how you want to handle the modifiers and the importance of preserving spacing, you might choose one method over the other. For this project, using Split allows us to preserve the original formatting more accurately, but we need to be careful in our processing logic to account for the empty strings.
	//         words := strings.Split(line, " ")
	//         */
	//         words := strings.Fields(line) // This will split the line into words while ignoring extra spaces.
	//         // if the line is empty, we should preserve it as is
	//         if len(words) == 0 {
	//             finalLines = append(finalLines, "")
	//             continue
	//         }
	//         // Process the modifiers in the words slice // Milestones 2, 3, 4,
	//         processedWords := processModifiers(words)
	//         // In The case we need to add Milestone 6 , Here the ORDER matters because we want to fix the punctuation spacing before we fix the quotes. If we fix the quotes first, we might end up with cases where a quote is attached to a punctuation mark, which could make it harder to correctly identify and format the quotes. By fixing the punctuation spacing first, we ensure that any punctuation is properly attached to the words before we handle the quotes, allowing for more accurate formatting of both punctuation and quotes in the final output.
	//         // Fix single quote formatting in the final words // Milestone 6
	//         quotedWords := fixQuotes(processedWords)
	//         // Fix punctuation spacing in the processed words // Milestone 5
	//         finalWords := fixPunctuationSpacing(quotedWords)
	//         // Join the words with spaces and save the line
	//         finalLines = append(finalLines, strings.Join(finalWords, " "))
	//     }
	// // Join the processed lines back into a single string with newlines
	// content = strings.Join(finalLines, "\n")

	// process all the transformations in a single function
	result := processText(content)

	// 4. Write the result string to the output file
	err = writeFile(outPutFile, result)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File processed successfully.")
}

// Helper to collect the text transormation functions
func processText(input string) string {
	lines := strings.Split(input, "\n")
	var finalLines []string // This will hold the processed lines

	for _, line := range lines {
		// split the line into words
		words := strings.Fields(line)
		// if the line is empty, we should preserve it as is
		if len(words) == 0 {
			finalLines = append(finalLines, "")
			continue
		}

		// Milestone 2, 3, 4: Modifiers
		words = processModifiers(words)
		// Milestone 6: Quotes
		words = fixQuotes(words)
		// Milestone 5: Punctuation (Splitting glued punct)
		words = fixPunctuationSpacing(words)
		// Milestone 7: Articles
		words = fixArticles(words)

		finalLines = append(finalLines, strings.Join(words, " "))
	}
	return strings.Join(finalLines, "\n")
}

// Milestone 1: Read Input, Write Output
func readFile(filename string) (string, error) {
	// Read the file
	date, err := os.ReadFile(filename)
	// Convert the result to a string and return it
	return string(date), err
}

func writeFile(filename string, content string) error {
	// Write the content to the file with appropriate permissions
	return os.WriteFile(filename, []byte(content), 0644)
}

// Milestone 2: Number Conversions: (hex) and (bin)
func hexToDecimal(s string) string {
	// Parse s as a base-16 integer
	result, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return s // If parsing fails, return s unchanged
	}
	// Convert the result back to a decimal string and return it
	return strconv.FormatInt(result, 10) // we use FormatInt to convert the integer back to a string in base 10
}

func binToDecimal(s string) string {
	// Parse s as a base-2 integer
	result, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s // If parsing fails, return s unchanged
	}
	// Convert the result back to a decimal string and return it
	return strconv.FormatInt(result, 10) // we use FormatInt to convert the integer back to a string in base 10
}

// Milestone 3: Case Modifiers: (up), (low), (cap) and Milestone 4: Numbered Modifiers: (up, 2), (cap, 4)
func processModifiers(words []string) []string {
	result := []string{}

	for i := 0; i < len(words); i++ {
		word := words[i]

		// 1. Handle Numbered Modifiers first (up, n), (low, n), (cap, n)
		// We check for the prefix because the word is "(up,"
		if (word == "(up," || word == "(low," || word == "(cap," || word == "(bin,") && i+1 < len(words) {
			if i+1 < len(words) {
				// Get the number and remove the closing ")" and any punctuation like "2)!"
				nStr := strings.Trim(words[i+1], ".,!?:;)")
				n, err := strconv.Atoi(nStr)

				if err == nil {
					for j := 1; j <= n; j++ {
						target := len(result) - j
						if target >= 0 {
							if word == "(up," {
								result[target] = strings.ToUpper(result[target])
							} else if word == "(low," {
								result[target] = strings.ToLower(result[target])
							} else if word == "(cap," {
								result[target] = capitalize(result[target])
							} else if word == "(bin," {
								result[target] = binToDecimal(result[target])
							}
						}
					}
				}
				i++ // Skip the number token
				continue
			}
		}

		// 2. Handle Simple Modifiers (hex), (bin), (up), (low), (cap)
		// Use Trim to handle "working (up)," -> "(up)"
		// First, strip any trailing punctuation so we can detect the modifier
		suffix := ""
		cleanWord := word
		for len(cleanWord) > 0 && strings.ContainsRune(".,!?:;", rune(cleanWord[len(cleanWord)-1])) {
			suffix = string(cleanWord[len(cleanWord)-1]) + suffix
			cleanWord = cleanWord[:len(cleanWord)-1]
		}

		switch cleanWord {
		case "(hex)":
			if len(result) > 0 {
				result[len(result)-1] = hexToDecimal(result[len(result)-1])
				result[len(result)-1] += suffix
			} else if suffix != "" {
				result = append(result, suffix)
			}
			continue
		case "(bin)":
			if len(result) > 0 {
				result[len(result)-1] = binToDecimal(result[len(result)-1])
				result[len(result)-1] += suffix
			} else if suffix != "" {
				result = append(result, suffix)
			}
			continue
		case "(up)":
			if len(result) > 0 {
				result[len(result)-1] = strings.ToUpper(result[len(result)-1])
				result[len(result)-1] += suffix
			} else if suffix != "" {
				result = append(result, suffix)
			}
			continue
		case "(low)":
			if len(result) > 0 {
				result[len(result)-1] = strings.ToLower(result[len(result)-1])
				result[len(result)-1] += suffix
			} else if suffix != "" {
				result = append(result, suffix)
			}
			continue
		case "(cap)":
			if len(result) > 0 {
				result[len(result)-1] = capitalize(result[len(result)-1])
				result[len(result)-1] += suffix
			} else if suffix != "" {
				result = append(result, suffix)
			}
			continue
		}

		// 3. If it's not a modifier, add it to result
		result = append(result, word)
	}
	return result
}

// Helper Function for Milestone 3 — Case Modifiers
func capitalize(s string) string {
	// If the string is empty, return it as is
	if len(s) == 0 {
		return s
	}
	lower := strings.ToLower(s)
	// Convert the first letter to uppercase and the rest to lowercase
	return strings.ToUpper(string(lower[0])) + lower[1:]
}

// Milestone 5 — Punctuation Spacing
func fixPunctuationSpacing(words []string) []string {
	result := []string{}
	puncs := ".,!?:;"

	for _, word := range words {
		if len(word) > 1 && strings.ContainsRune(puncs, rune(word[0])) && !isPunctuation(word) { // if the word starts with a punctuation mark but is not entirely punctuation, we need to split it
			pEnd := 0 // pointer to find the end of the punctuation prefix
			// Move pEnd until we find a character that is not a punctuation mark
			for pEnd < len(word) && strings.ContainsRune(puncs, rune(word[pEnd])) { // Check if the character at pEnd is a punctuation mark
				pEnd++
			}

			prefix := word[:pEnd] // The punctuation prefix that we need to attach to the previous word
			rest := word[pEnd:]   // The rest of the word that should be added as a new element in the result slice

			if len(result) > 0 {
				result[len(result)-1] += prefix
			} else {
				result = append(result, prefix)
			}

			result = append(result, rest)

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

// Helper function for Milestone 5
func isPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	// only return true if the entire string consists of allowed marks
	for _, r := range s {
		if !strings.ContainsRune(".,!?:;", r) { // strings.ContainsRune is a function that checks if a rune is in a string and returns true if it is, false otherwise
			return false
		}
	}
	return true
}

// Milestone 6 — Single Quote Formatting
func fixQuotes(words []string) []string {
	result := []string{}

	quoteOpen := false

	for i := 0; i < len(words); i++ {
		word := words[i]

		if word == "'" {
			if !quoteOpen {
				// Opening quote: glue it to the next word
				if i+1 < len(words) {
					words[i+1] = "'" + words[i+1]
					quoteOpen = true
					continue // Skip adding the standalone quote to result
				}
			} else {
				// Closing quote: glue it to the previous word
				if len(result) > 0 {
					result[len(result)-1] += "'"
					quoteOpen = false
					continue // Skip adding the standalone quote to result
				}
			}
		}

		if strings.HasPrefix(word, "'") && len(word) > 1 {
			quoteOpen = true
		} // If the word starts with a quote, we consider it an opening quote. This handles cases where the quote is already attached to a word, like "'hello". We check if the length of the word is greater than 1 to ensure that it's not just a standalone quote.
		if strings.HasSuffix(word, "'") && len(word) > 1 {
			quoteOpen = false
		} // If the word ends with a quote, we consider it a closing quote. This handles cases where the quote is attached to the end of a word, like "world'". Again, we check if the length of the word is greater than 1 to ensure that it's not just a standalone quote.

		// If it's not a standalone quote, add it to result
		result = append(result, word)

	}
	return result
}

// Milestone 7: Article Correction (A -> An)
func fixArticles(words []string) []string {
	vowels := "aeiouhAEIOUH"
	for i := 0; i < len(words)-1; i++ {
		// Handle both 'a' and 'A' even if wrapped in a quote (like "'a")
		clean := strings.Trim(words[i], "'\"")
		prefix := words[i][:strings.Index(words[i], clean)]

		if strings.ToLower(clean) == "a" {
			next := words[i+1]
			// Find the first actual letter in the next word (skip quotes)
			letterIdx := 0
			for letterIdx < len(next) && !((next[letterIdx] >= 'a' && next[letterIdx] <= 'z') || (next[letterIdx] >= 'A' && next[letterIdx] <= 'Z')) {
				letterIdx++
			}

			if letterIdx < len(next) && strings.ContainsRune(vowels, rune(next[letterIdx])) {
				if clean == "a" {
					words[i] = prefix + "an"
				} else {
					words[i] = prefix + "An"
				}
			}
		}
	}
	return words
}
