package varName

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// TextAlignment refers to which section of text we want to prioritize
type TextAlignment int

// NameCase refers to one of a set of casing schemes such as camel, snake, kebab
type NameCase int

const (
	// Left keeps left-most text and truncates anything after a specified limit is exceeded
	Left TextAlignment = iota
	// Right keeps right-most text and truncates anything after a specified limit is exceeded
	Right
	// Edge keeps outer text and truncates any interior text after a specified limit is exceeded
	Edge
)

const (
	// Camel refers to camel case; eg. somethingCamelCased
	Camel NameCase = iota
	// Snake refers to snake case; eg. something_snake_cased
	Snake
	// Kebab refers to kebab case; eg. something-kebab-cased
	Kebab
)

var (
	// matches substrings containing 2 or more consecutive capital letters
	multiCapsPattern = regexp.MustCompile(`([A-Z]{2,})`)
	// matches single capital letters
	capsPattern = regexp.MustCompile(`([A-Z])`)
	// matches substrings containing 2 or more consecutive spaces
	multiSpacePattern = regexp.MustCompile(`\s{2,}`)
	// stop words/ skip words to filter out
	defaultSkipwords = map[string]bool{
		// A,
		"a": true, "about": true, "above": true, "after": true, "again": true, "against": true, "all": true, "am": true, "an": true, "and": true, "any": true, "are": true, "aren't": true, "as": true, "at": true,
		// B, C,
		"be": true, "because": true, "been": true, "before": true, "being": true, "below": true, "between": true, "both": true, "but": true, "by": true, "can't": true, "cannot": true, "could": true, "couldn't": true,
		// D, E, F, G,
		"did": true, "didn't": true, "do": true, "does": true, "doesn't": true, "doing": true, "don't": true, "down": true, "during": true, "each": true, "few": true, "for": true, "from": true, "further": true,
		// H, I, J, K
		"had": true, "hadn't": true, "has": true, "hasn't": true, "have": true, "haven't": true, "having": true, "he": true, "he'd": true, "he'll": true, "he's": true, "her": true, "here": true, "here's": true, "hers": true, "herself": true, "him": true, "himself": true, "his": true, "how": true, "how's": true, "i'd": true, "i'll": true, "i'm": true, "i've": true, "if": true, "in": true, "into": true, "is": true, "isn't": true, "it": true, "it's": true, "its": true, "itself": true,
		// L, M, N, O, P, Q, R
		"let's": true, "more": true, "most": true, "mustn't": true, "no": true, "nor": true, "not": true, "of": true, "off": true, "on": true, "once": true, "only": true, "or": true, "other": true, "ought": true, "our": true, "ours": true, "ourselves": true, "out": true, "over": true, "own": true,
		// S, T,
		"same": true, "shan't": true, "she": true, "she'd": true, "she'll": true, "she's": true, "should": true, "shouldn't": true, "so": true, "some": true, "such": true, "than": true, "that": true, "that's": true, "the": true, "their": true, "theirs": true, "them": true, "themselves": true, "then": true, "there": true, "there's": true, "these": true, "they": true, "they'd": true, "they'll": true, "they're": true, "they've": true, "this": true, "those": true, "through": true, "to": true, "too": true,
		// U, V, W, X, Y, Z
		"under": true, "until": true, "up": true, "very": true, "was": true, "wasn't": true, "we'd": true, "we'll": true, "we're": true, "we've": true, "were": true, "weren't": true, "what": true, "what's": true, "when": true, "when's": true, "where": true, "where's": true, "which": true, "while": true, "who": true, "who's": true, "whom": true, "why": true, "why's": true, "with": true, "won't": true, "would": true, "wouldn't": true, "you": true, "you'd": true, "you'll": true, "you're": true, "you've": true, "your": true, "yours": true, "yourself": true, "yourselves": true,
	}

	defaultSubstitutions = map[string]string{
		// special backslash
		`\`: "",
		// multi-symbol equivalence operators
		`<>`: "ne", `<=`: "lte", `>=`: "gte", `~=`: "ne", `!=`: "ne", `^=`: "ne",
		// single-symbol equivalence operators
		`=`: "eq", `<`: "lt", `>`: "gt",
		// other interpretted symbols
		`%`: "pct", `&`: "and", `/`: "per", `US$`: "usd", `$`: "usd",
		// single symbols, replace with space
		`-`: " ", `_`: " ",
		// single symvols, replace with ""
		`.`: "", `#`: "", `?`: "", `|`: "", `*`: "", `,`: "", `(`: "", `)`: "", `:`: "", `;`: "", `'`: "", `"`: "",
	}

	numberedNamePattern   = regexp.MustCompile(`(.*)_(\d+)$`)
	numberedSuffixPattern = regexp.MustCompile(`_\d+$`)
)

// String Cleaning functions
func convertMultiCapsToSingleCaps(s string) string {
	multiCapsToSingleCapWithSpace := func(s string) string {
		return fmt.Sprintf("%s ", strings.Title(strings.ToLower(s)))
	}
	return multiCapsPattern.ReplaceAllStringFunc(s, multiCapsToSingleCapWithSpace)
}

func interpretCamelCaseAsSpace(s string) string {
	insertSpaceBefore := func(s string) string { return fmt.Sprintf(" %s", s) }
	return capsPattern.ReplaceAllStringFunc(s, insertSpaceBefore)
}

func removeMultiSpaceAndTrim(s string) string {
	return multiSpacePattern.ReplaceAllString(strings.TrimSpace(s), " ")
}

func isNumOrSpace(c rune) bool {
	return unicode.IsSpace(c) || unicode.IsNumber(c)
}

// parseExistingCamelDelim takes a string
// such as     "EconIndicatorNominalGDP1997China"
// and outputs "Econ Indicator Nominal Gdp 1997 China"
func parseExistingCamelDelim(s string) string {
	// convert multicaps to single caps
	s = convertMultiCapsToSingleCaps(s)
	// interpret camelCase as space
	s = interpretCamelCaseAsSpace(s)
	// remove duplicate spaces and trim
	s = removeMultiSpaceAndTrim(s)
	return s
}

// remapChars iterates through a map of stubstitutions replacing any present characters in s
// accordingly.  the removeOnly option allows you to override the mapping to remap to an empty string ""
func remapChars(s string, substitutions *map[string]string, removeOnly bool) string {
	for origVal, replVal := range *substitutions {
		if removeOnly {
			replVal = ""
		}
		s = strings.Replace(s, origVal, replVal, -1)
	}
	s = removeMultiSpaceAndTrim(s)
	return s
}

func reverse(ss []string) []string {
	last := len(ss) - 1
	reversed := []string{}
	for i := last; i >= 0; i-- {
		reversed = append(reversed, ss[i])
	}
	return reversed
}

// truncateList takes a slice of strings and allows you to truncate
// either the left fright or middle of the list to satisfy a maximum combined
// string length, maxLen.  Function currently accepts parameters indidating
// which portion of the list you want to *keep* (Left to keep left, truncate right,
// Edge to truncate middle etc.)
// TODO: should revisit this, if confusing, can refactor to spcify what you want to *cut*
func truncateList(wordList []string, maxLen int, alignment TextAlignment) []string {
	// listLen := len(wordList)
	reversedList := reverse(wordList)
	truncatedList := []string{}
	//tailList := []string{}
	charCount := 0
	if alignment == Edge {
		head := []string{}
		tail := []string{}
		for i := range wordList {
			h := wordList[i]
			t := reversedList[i]
			if len(h)+len(t)+charCount <= maxLen {
				head = append(head, h)
				tail = append(tail, t)
				charCount += len(h) + len(t)
			} else if len(h)+charCount <= maxLen {
				head = append(head, h)
				charCount += len(h)
			} else {
				break
			}
		}
		truncatedList = append(head, reverse(tail)...)
	} else {
		if alignment == Right {
			wordList = reversedList
		}
		for _, w := range wordList {
			if len(w)+charCount > maxLen {
				break
			}
			truncatedList = append(truncatedList, w)
			charCount += len(w)
		}
		if alignment == Right {
			truncatedList = reverse(truncatedList)
		}
	}
	return truncatedList
}

// listToVarName processes a slice of words and rejoins them into a single string appropriate for a single-word variable name
func listToVarName(wordList []string, skipwords *map[string]bool, maxLen int, alignment TextAlignment, noRepeats bool, caseType NameCase) string {
	// ensure skipwords are lower-cased
	skipwordsLower := map[string]bool{}
	for key := range *skipwords {
		lKey := strings.ToLower(key)
		skipwordsLower[lKey] = true
	}
	// filter wordList
	filteredList := []string{}
	for _, word := range wordList {
		wordL := strings.ToLower(word)
		shouldSkip := skipwordsLower[wordL]
		if len(wordL) > 0 && !shouldSkip {
			filteredList = append(filteredList, strings.ToTitle(wordL))
			//if noRepats == true, prevent repeated words
			if noRepeats {
				skipwordsLower[wordL] = true
			}
		}
	}
	// truncate filteredList
	truncatedList := truncateList(wordList, maxLen, alignment)
	// join truncatedlist
	switch caseType {
	case Camel:
		//make camel case
		camelCaseList := []string{}
		for i, word := range truncatedList {
			if i == 0 {
				camelCaseList = append(camelCaseList, strings.ToLower(word))
			} else {
				camelCaseList = append(camelCaseList, strings.Title(word))
			}
		}
		return strings.Join(camelCaseList, "")
	case Kebab:
		return strings.Trim(strings.ToLower(strings.Join(truncatedList, "-")), "- ")
	default: // (Snake)
		return strings.Trim(strings.ToLower(strings.Join(truncatedList, "_")), "_ ")
	}

}

// NameParams is a convenience parameter struct made for CreateVarNameFromParams
type NameParams struct {
	InputName string
	// dictionary of strings to filter out
	SkipWords *map[string]bool
	// dictioanry of strings to be substitute
	Substitutions *map[string]string
	// delimiter to use as word boundary
	Delim string
	// maximum lenth used
	MaxLen int
	// the if you want to override the replacement values with "", set to true
	RemoveOnly bool
	// if we should not allow words to occur repeated times, set to true
	NoRepeats bool
	// a choice of Left, Right, or Edge specifies which region of the input to prioritize
	Alignment TextAlignment
	// a choice of Camel, Snake, or Kebab to determine the name-casing
	NameCasing NameCase
}

// NewNameParams creates a pointer to a new NameParams struct and
// populates it with a set of default options
func NewNameParams(name string) *NameParams {
	return &NameParams{
		InputName:     name,
		SkipWords:     &defaultSkipwords,
		Substitutions: &defaultSubstitutions,
		Delim:         " ",
		MaxLen:        30,
		RemoveOnly:    false,
		NoRepeats:     true,
		Alignment:     Left,
		NameCasing:    Snake,
	}
}

// CreateVarNameFromParams takes a lengthy title string and attempts to generate a
// condensed but still recognizable variable name
func CreateVarNameFromParams(p *NameParams) string {
	s := p.InputName
	s = parseExistingCamelDelim(s)
	s = remapChars(s, p.Substitutions, p.RemoveOnly)
	// split current string on Delim
	wordList := strings.Split(s, p.Delim)
	// filter, truncate, and rejoin
	tableName := listToVarName(wordList, p.SkipWords, p.MaxLen, p.Alignment, p.NoRepeats, p.NameCasing)
	// remove any leading numbers
	tableName = strings.TrimLeftFunc(tableName, isNumOrSpace)
	// remove any extraneous underscores or dashes
	tableName = strings.Trim(tableName, "_- ")
	return tableName
}

// CreateVarNameFromString does the same as CreateVarNameFromParams but
// uses a set of default settings and only reuires a string as an input
func CreateVarNameFromString(s string) string {
	p := NewNameParams(s)
	return CreateVarNameFromParams(p)
}

// MakeNameUnique takes a name, evaluates it against a set of existing names and
// returns a unique name by appending an underscore and number to the end
// Each time the function generates a new name, the name is added to the existing map
func MakeNameUnique(name string, existing *map[string]bool) string {
	// break name into base and suffix
	baseName := numberedNamePattern.ReplaceAllString(name, "$1")
	isUnique := true
	if (*existing)[name] || (*existing)[baseName] {
		isUnique = false
	}
	highestNum := 1
	for existingName := range *existing {
		matches := numberedNamePattern.FindStringSubmatch(existingName)

		// Note we expect a slice of len=3 b/c there are 2 capturing groups in
		// this pattern -- /(.*)_(\d+)$/ -- abc_123 -> [abc_123, abc, 123]
		if len(matches) == 3 {
			if matches[1] == baseName {
				isUnique = false
				numStr := matches[2]
				// convert to int if possible
				i, err := strconv.Atoi(numStr)
				if err == nil {
					// if higher then highestNum, update
					if i > highestNum {
						highestNum = i
					}
				}
			}
		}
	}
	uniqueName := fmt.Sprintf("%s_%d", baseName, highestNum+1)
	if isUnique {
		uniqueName = name
	}
	(*existing)[uniqueName] = true
	return uniqueName
}
