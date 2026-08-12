package engine

import (
	"fmt"
	"regexp"
	"strings"
)

// natvaRegex converts 'n' -> 'R' after triggering consonants (r, f, F, z)
// Note: Go's regexp package does not support lookaheads (?=...), so we capture
// the following vowel/consonant in Group 2 and replace with ${1}R${2}.
var natvaRegex = regexp.MustCompile(`([rfFz][aAiIuUfFxXeEoOkKgGNpPbBmyvhM]*)n([aAiIuUfFxXeEoOmyv])`)

func applyNatva(word string) string {
	current := word
	for natvaRegex.MatchString(current) {
		current = natvaRegex.ReplaceAllString(current, "${1}R${2}")
	}
	return current
}

func buildGrid(pr, dv, tr, ca, pa, sa, sap, sam [3]string) map[string][]string {
	return map[string][]string{
		"prathama":   {pr[0], pr[1], pr[2]},
		"dvitiya":    {dv[0], dv[1], dv[2]},
		"tritiya":    {tr[0], tr[1], tr[2]},
		"caturthi":   {ca[0], ca[1], ca[2]},
		"panchami":   {pa[0], pa[1], pa[2]},
		"sasthi":     {sa[0], sa[1], sa[2]},
		"saptami":    {sap[0], sap[1], sap[2]},
		"sambodhana": {sam[0], sam[1], sam[2]},
	}
}

// DeclineNoun generates the 8x3 declension grid for a nominal stem
func DeclineNoun(base string, gender string) (map[string][]string, error) {
	if irreg, ok := declineIrregular(base, gender); ok {
		fixed := make(map[string][]string)
		for k, v := range irreg {
			natvaV := make([]string, len(v))
			for i, w := range v {
				natvaV[i] = applyNatva(w)
			}
			fixed[k] = natvaV
		}
		return fixed, nil
	}

	res := make(map[string][]string)
	s := ""
	if len(base) > 0 {
		s = base[:len(base)-1]
	}

	if strings.HasSuffix(base, "a") {
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s + "aH", s + "O", s + "AH"},
				[3]string{s + "am", s + "O", s + "An"},
				[3]string{s + "ena", s + "AByAm", s + "EH"},
				[3]string{s + "Aya", s + "AByAm", s + "eByaH"},
				[3]string{s + "At", s + "AByAm", s + "eByaH"},
				[3]string{s + "asya", s + "ayoH", s + "AnAm"},
				[3]string{s + "e", s + "ayoH", s + "ezu"},
				[3]string{s + "a", s + "O", s + "AH"},
			)
		} else if gender == "neuter" {
			res = buildGrid(
				[3]string{s + "am", s + "e", s + "Ani"},
				[3]string{s + "am", s + "e", s + "Ani"},
				[3]string{s + "ena", s + "AByAm", s + "EH"},
				[3]string{s + "Aya", s + "AByAm", s + "eByaH"},
				[3]string{s + "At", s + "AByAm", s + "eByaH"},
				[3]string{s + "asya", s + "ayoH", s + "AnAm"},
				[3]string{s + "e", s + "ayoH", s + "ezu"},
				[3]string{s + "a", s + "e", s + "Ani"},
			)
		}
	} else if strings.HasSuffix(base, "A") {
		if gender == "feminine" {
			res = buildGrid(
				[3]string{s + "A", s + "e", s + "AH"},
				[3]string{s + "Am", s + "e", s + "AH"},
				[3]string{s + "ayA", s + "AByAm", s + "ABiH"},
				[3]string{s + "AyE", s + "AByAm", s + "AByaH"},
				[3]string{s + "AyAH", s + "AByAm", s + "AByaH"},
				[3]string{s + "AyAH", s + "ayoH", s + "AnAm"},
				[3]string{s + "AyAm", s + "ayoH", s + "Asu"},
				[3]string{s + "e", s + "e", s + "AH"},
			)
		}
	} else if strings.HasSuffix(base, "I") {
		if gender == "feminine" {
			res = buildGrid(
				[3]string{s + "I", s + "yO", s + "yaH"},
				[3]string{s + "Im", s + "yO", s + "IH"},
				[3]string{s + "yA", s + "IByAm", s + "IBiH"},
				[3]string{s + "yE", s + "IByAm", s + "IByaH"},
				[3]string{s + "yAH", s + "IByAm", s + "IByaH"},
				[3]string{s + "yAH", s + "yoH", s + "InAm"},
				[3]string{s + "yAm", s + "yoH", s + "Izu"},
				[3]string{s + "i", s + "yO", s + "yaH"},
			)
		}
	} else if strings.HasSuffix(base, "i") {
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s + "iH", s + "I", s + "ayaH"},
				[3]string{s + "im", s + "I", s + "In"},
				[3]string{s + "inA", s + "iByAm", s + "iBiH"},
				[3]string{s + "aye", s + "iByAm", s + "iByaH"},
				[3]string{s + "eH", s + "iByAm", s + "iByaH"},
				[3]string{s + "eH", s + "yoH", s + "InAm"},
				[3]string{s + "O", s + "yoH", s + "izu"},
				[3]string{s + "e", s + "I", s + "ayaH"},
			)
		} else if gender == "feminine" {
			res = buildGrid(
				[3]string{s + "iH", s + "I", s + "ayaH"},
				[3]string{s + "im", s + "I", s + "IH"},
				[3]string{s + "yA", s + "iByAm", s + "iBiH"},
				[3]string{s + "yE", s + "iByAm", s + "iByaH"},
				[3]string{s + "yAH", s + "iByAm", s + "iByaH"},
				[3]string{s + "yAH", s + "yoH", s + "InAm"},
				[3]string{s + "yAm", s + "yoH", s + "izu"},
				[3]string{s + "e", s + "I", s + "ayaH"},
			)
		}
	} else if strings.HasSuffix(base, "at") {
		s2 := ""
		if len(base) >= 2 {
			s2 = base[:len(base)-2]
		}
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s2 + "an", s2 + "antO", s2 + "antaH"},
				[3]string{s2 + "antam", s2 + "antO", s2 + "ataH"},
				[3]string{s2 + "atA", s2 + "adByAm", s2 + "adBiH"},
				[3]string{s2 + "ate", s2 + "adByAm", s2 + "adByaH"},
				[3]string{s2 + "ataH", s2 + "adByAm", s2 + "adByaH"},
				[3]string{s2 + "ataH", s2 + "atoH", s2 + "atAm"},
				[3]string{s2 + "ati", s2 + "atoH", s2 + "atsu"},
				[3]string{s2 + "an", s2 + "antO", s2 + "antaH"},
			)
		}
	} else if strings.HasSuffix(base, "an") {
		s2 := ""
		if len(base) >= 2 {
			s2 = base[:len(base)-2]
		}
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s2 + "A", s2 + "AnO", s2 + "AnaH"},
				[3]string{s2 + "Anam", s2 + "AnO", s2 + "naH"},
				[3]string{s2 + "nA", s2 + "aByAm", s2 + "aBiH"},
				[3]string{s2 + "ne", s2 + "aByAm", s2 + "aByaH"},
				[3]string{s2 + "naH", s2 + "aByAm", s2 + "aByaH"},
				[3]string{s2 + "naH", s2 + "noH", s2 + "nAm"},
				[3]string{s2 + "ni", s2 + "noH", s2 + "asu"},
				[3]string{s2 + "an", s2 + "AnO", s2 + "AnaH"},
			)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("declension logic for '%s' in gender '%s' is not implemented", base, gender)
	}

	for k, forms := range res {
		for i, word := range forms {
			forms[i] = applyNatva(word)
		}
		res[k] = forms
	}

	return res, nil
}