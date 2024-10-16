package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func trimNonAlphanumeric(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-'
	})
}

/*
To sort a slice will use the following struct.
*/
type WordCount struct {
	word  string
	count int
}

/*
This function checks if givern word looks like complete word. Basically it mostly needed
for checking for the requirement '* "-" словом не является'.
Generalize this requirement to check that any single non-alpanumeric character is not a word.
*/
func isCompleteWord(word string) bool {
	var r rune
	for _, c := range word {
		r = c
		break
	}
	if len(word) == 1 && !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
		return false
	}
	return true
}

/*
For better testability we add this intermediate function, which returns a sorted slice of pairs
word->count packed in WordCount structure.
*/
func WordsWithCount(input string) []WordCount {
	// Split input text by words by whitespace characters
	pieces := strings.Fields(input)

	// Prepare a map to count occurrences for words
	counter := make(map[string]int)

	for _, str := range pieces {
		//  remove all non-alphanumeric charaters from both word sizes
		word := trimNonAlphanumeric(str)

		// Check if it's really a word what's left
		if !isCompleteWord(word) {
			continue
		}

		// Satisfy the requirement 'не учитывать регистр букв'
		word = strings.ToLower(word)

		if word != "" {
			counter[word]++
		}
	}

	// Prepare a slice where we can sort the words. Reserved is already known for the slice
	words := make([]WordCount, 0, len(counter))

	for k, v := range counter {
		words = append(words, WordCount{k, v})
	}

	// Sort the slice based on the Value field in descending order
	sort.Slice(words, func(x, y int) bool {
		if words[x].count == words[y].count {
			return words[x].word < words[y].word
		}
		return words[x].count > words[y].count
	})

	return words
}

/*
Return top 10 most often used words from given text.
*/
func Top10(input string) []string {
	words := WordsWithCount(input)
	// Prepare an output slice of words which is not more than 10 elements
	nElements := min(len(words), 10)
	ret := make([]string, 0, nElements)

	for i := 0; i < nElements; i++ {
		ret = append(ret, words[i].word)
	}

	return ret
}
