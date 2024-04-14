package words

import (
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
)

var Contractions = map[string]string{
	"aren't":    "are not",
	"can't":     "cannot",
	"couldn't":  "could not",
	"didn't":    "did not",
	"doesn't":   "does not",
	"don't":     "do not",
	"hadn't":    "had not",
	"hasn't":    "has not",
	"haven't":   "have not",
	"he'd":      "he had",
	"he'll":     "he will",
	"he's":      "he is",
	"i'd":       "i had",
	"i'll":      "i will",
	"i'm":       "i am",
	"i've":      "i have",
	"isn't":     "is not",
	"it'd":      "it had",
	"it'll":     "it will",
	"it's":      "it is",
	"let's":     "let us",
	"mightn't":  "might not",
	"mustn't":   "must not",
	"shan't":    "shall not",
	"she'd":     "she had",
	"she'll":    "she will",
	"she's":     "she is",
	"shouldn't": "should not",
	"that's":    "that is",
	"there's":   "there is",
	"they'd":    "they had",
	"they'll":   "they will",
	"they're":   "they are",
	"they've":   "they have",
	"wasn't":    "was not",
	"we'd":      "we had",
	"we'll":     "we will",
	"we're":     "we are",
	"we've":     "we have",
	"weren't":   "were not",
	"what'll":   "what will",
	"what're":   "what are",
	"what's":    "what is",
	"what've":   "what have",
	"where's":   "where is",
	"who'd":     "who had",
	"who'll":    "who will",
	"who's":     "who is",
	"who've":    "who have",
	"won't":     "will not",
	"wouldn't":  "would not",
	"you'd":     "you had",
	"you'll":    "you will",
	"you're":    "you are",
	"you've":    "you have",
}

func ExpandContractions(word string) string {
	if expanded, exists := Contractions[word]; exists {
		return expanded
	}
	return word
}

func SplitString(input string) []string {
	reg := regexp.MustCompile("[^a-zA-Z']+")
	return reg.Split(input, -1)

}

func Stemming(words []string) ([]string, error) {
	var normalized []string
	contains := make(map[string]bool)
	for _, v := range words {
		v = strings.ToLower(v)
		newLine := ExpandContractions(v)
		expandedWords := strings.Fields(newLine)
		for _, word := range expandedWords {
			if !english.IsStopWord(word) {
				stemmed, err := snowball.Stem(word, "english", false)
				if err != nil {
					return nil, err
				}
				if stemmed != "" {
					if _, ok := contains[stemmed]; !ok {
						normalized = append(normalized, stemmed)
						contains[stemmed] = true
					}
				}
			}
		}
	}
	return normalized, nil
}
