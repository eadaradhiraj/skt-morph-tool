package engine

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type upasargaRule struct {
	pattern   *regexp.Regexp
	canonical string
	prepend   string
}

var upasargaPatterns []upasargaRule

func init() {
	upasargaPatterns = []upasargaRule{
		{pattern: regexp.MustCompile(`^samudA`), canonical: "sam + ud + AN", prepend: ""},
		{pattern: regexp.MustCompile(`^sa[mMnYNR]u[dtcjNYRnl]A`), canonical: "sam + ud + AN", prepend: ""},
		{pattern: regexp.MustCompile(`^sa[mMnYNR]u[dtcjNYRnl]`), canonical: "sam + ud", prepend: ""},
		{pattern: regexp.MustCompile(`^vyA`), canonical: "vi + AN", prepend: ""},
		{pattern: regexp.MustCompile(`^pratyA`), canonical: "prati + AN", prepend: ""},
		{pattern: regexp.MustCompile(`^vy`), canonical: "vi", prepend: ""},
		{pattern: regexp.MustCompile(`^vi`), canonical: "vi", prepend: ""},
		{pattern: regexp.MustCompile(`^praty`), canonical: "prati", prepend: ""},
		{pattern: regexp.MustCompile(`^prati`), canonical: "prati", prepend: ""},
		{pattern: regexp.MustCompile(`^sa[mMnYNR]`), canonical: "sam", prepend: ""},
		{pattern: regexp.MustCompile(`^sam`), canonical: "sam", prepend: ""},
		{pattern: regexp.MustCompile(`^pra`), canonical: "pra", prepend: ""},
		{pattern: regexp.MustCompile(`^A`), canonical: "AN", prepend: ""},
	}
}

// GetUpasargaSplits returns potential upasarga splits as pairs of [canonical, stripped_stem]
func GetUpasargaSplits(word string) [][2]string {
	var splits [][2]string
	for _, r := range upasargaPatterns {
		loc := r.pattern.FindStringIndex(word)
		if loc != nil {
			endIdx := loc[1]
			stripped := r.prepend + word[endIdx:]
			if utf8.RuneCountInString(stripped) >= 2 {
				splits = append(splits, [2]string{r.canonical, stripped})
			}
		}
	}
	return splits
}

// ApplyUpasargaSandhi applies sandhi rules when combining prefixes with verb forms
func ApplyUpasargaSandhi(upasargaStr string, form string) string {
	if upasargaStr == "" {
		return form
	}

	rawPrefixes := strings.Split(upasargaStr, "+")
	prefixes := make([]string, 0, len(rawPrefixes))
	for _, s := range rawPrefixes {
		prefixes = append(prefixes, strings.TrimSpace(s))
	}

	result := form

	// Iterate prefixes in reverse (matching Rust's rev())
	for i := len(prefixes) - 1; i >= 0; i-- {
		p := prefixes[i]
		if p == "AN" {
			p = "A"
		}

		pRunes := []rune(p)
		if len(pRunes) == 0 {
			continue
		}

		pLast := pRunes[len(pRunes)-1]

		rRunes := []rune(result)
		var rFirst rune
		if len(rRunes) > 0 {
			rFirst = rRunes[0]
		}

		isVowel := func(c rune) bool {
			return strings.ContainsRune("aAiIuUfFeEoO", c)
		}

		pMinusOne := string(pRunes[:len(pRunes)-1])
		rMinusFirst := ""
		if len(rRunes) > 1 {
			rMinusFirst = string(rRunes[1:])
		}

		if pLast == 'i' && isVowel(rFirst) {
			result = pMinusOne + "y" + result
		} else if pLast == 'u' && isVowel(rFirst) {
			result = pMinusOne + "v" + result
		} else if pLast == 'a' || pLast == 'A' {
			if rFirst == 'a' || rFirst == 'A' {
				result = pMinusOne + "A" + rMinusFirst
			} else if rFirst == 'i' || rFirst == 'I' {
				result = pMinusOne + "e" + rMinusFirst
			} else if rFirst == 'u' || rFirst == 'U' {
				result = pMinusOne + "o" + rMinusFirst
			} else if rFirst == 'f' || rFirst == 'F' {
				result = pMinusOne + "ar" + rMinusFirst
			} else {
				result = p + result
			}
		} else if pLast == 'm' {
			if !isVowel(rFirst) {
				result = pMinusOne + "M" + result
			} else {
				result = p + result
			}
		} else if pLast == 'd' {
			if strings.ContainsRune("kKqQpPzSs", rFirst) {
				result = pMinusOne + "t" + result
			} else if strings.ContainsRune("cC", rFirst) {
				result = pMinusOne + "c" + result
			} else if strings.ContainsRune("jJ", rFirst) {
				result = pMinusOne + "j" + result
			} else {
				result = p + result
			}
		} else {
			result = p + result
		}

		// applyNatva is available in package scope from declension.go
		result = applyNatva(result)
	}

	return result
}