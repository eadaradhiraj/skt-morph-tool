package engine

import (
	"fmt"
	"strings"
)

// applyNatva applies Sanskrit Natva (n -> R) rule according to Paninian grammar (8.4.1 - 8.4.2)
// Triggers: r, f, F, z
// Intervenors allowed: Vowels, Gutturals (k, K, g, G, N), Labials (p, P, b, B, m), Glides/Aspirates (y, v, h, M)
// Target: n followed by vowel/glide/nasal
func applyNatva(word string) string {
	runes := []rune(word)
	n := len(runes)
	if n == 0 {
		return word
	}

	isTrigger := func(r rune) bool {
		return r == 'r' || r == 'f' || r == 'F' || r == 'z'
	}

	isAllowedIntervenor := func(r rune) bool {
		if strings.ContainsRune("aAiIuUfFxXeEoO", r) {
			return true
		}
		if strings.ContainsRune("kKgGN", r) {
			return true
		}
		if strings.ContainsRune("pPbBm", r) {
			return true
		}
		if strings.ContainsRune("yvhM", r) {
			return true
		}
		return false
	}

	isFollowingEnv := func(r rune) bool {
		return strings.ContainsRune("aAiIuUfFxXeEoOmyvn", r)
	}

	result := make([]rune, n)
	copy(result, runes)
	hasTrigger := false

	for i := 0; i < n; i++ {
		r := runes[i]
		if isTrigger(r) {
			hasTrigger = true
			continue
		}

		if hasTrigger {
			if r == 'n' {
				if i+1 < n && isFollowingEnv(runes[i+1]) {
					result[i] = 'R'
				}
				hasTrigger = false
			} else if !isAllowedIntervenor(r) {
				hasTrigger = false
			}
		}
	}

	return string(result)
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

func declineIrregular(base string, gender string) (map[string][]string, bool) {
	if base == "go" && (gender == "masculine" || gender == "feminine") {
		return buildGrid(
			[3]string{"gOH", "gAvO", "gAvaH"},
			[3]string{"gAm", "gAvO", "gAH"},
			[3]string{"gavA", "goByAm", "goBiH"},
			[3]string{"gave", "goByAm", "goByaH"},
			[3]string{"goH", "goByAm", "goByaH"},
			[3]string{"goH", "gavoH", "gavAm"},
			[3]string{"gavi", "gavoH", "gozu"},
			[3]string{"gOH", "gAvO", "gAvaH"},
		), true
	}
	if base == "strI" && gender == "feminine" {
		return buildGrid(
			[3]string{"strI", "striyO", "striyaH"},
			[3]string{"striyam", "striyO", "striyaH"},
			[3]string{"striyA", "strIByAm", "strIBiH"},
			[3]string{"striyE", "strIByAm", "strIByaH"},
			[3]string{"striyAH", "strIByAm", "strIByaH"},
			[3]string{"striyAH", "striyoH", "strInAm"},
			[3]string{"striyAm", "striyoH", "strIzu"},
			[3]string{"stri", "striyO", "striyaH"},
		), true
	}
	if base == "sakhi" && gender == "masculine" {
		return buildGrid(
			[3]string{"saKA", "saKAnO", "saKAnaH"},
			[3]string{"saKAnam", "saKAnO", "saKIn"},
			[3]string{"saKyA", "saKiByAm", "saKiBiH"},
			[3]string{"saKye", "saKiByAm", "saKiByaH"},
			[3]string{"saKyuH", "saKiByAm", "saKiByaH"},
			[3]string{"saKyuH", "saKyoH", "saKInAm"},
			[3]string{"saKyO", "saKyoH", "saKizu"},
			[3]string{"saKe", "saKAnO", "saKAnaH"},
		), true
	}
	if base == "pati" && gender == "masculine" {
		return buildGrid(
			[3]string{"patiH", "patI", "patayaH"},
			[3]string{"patim", "patI", "patIn"},
			[3]string{"patyA", "patiByAm", "patiBiH"},
			[3]string{"patye", "patiByAm", "patiByaH"},
			[3]string{"patyuH", "patiByAm", "patiByaH"},
			[3]string{"patyuH", "patyoH", "patInAm"},
			[3]string{"patyO", "patyoH", "patizu"},
			[3]string{"pate", "patI", "patayaH"},
		), true
	}
	return nil, false
}

// DeclineNoun generates the 8x3 declension grid for a nominal stem
func DeclineNoun(base string, gender string) (map[string][]string, error) {
	if base == "" {
		return nil, fmt.Errorf("base stem cannot be empty")
	}

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

	if strings.HasSuffix(base, "a") {
		s := base[:len(base)-1]
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
		s := base[:len(base)-1]
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
		s := base[:len(base)-1]
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
		s := base[:len(base)-1]
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
		} else if gender == "neuter" {
			res = buildGrid(
				[3]string{s + "i", s + "inI", s + "Ini"},
				[3]string{s + "i", s + "inI", s + "Ini"},
				[3]string{s + "inA", s + "iByAm", s + "iBiH"},
				[3]string{s + "ine", s + "iByAm", s + "iByaH"},
				[3]string{s + "inaH", s + "iByAm", s + "iByaH"},
				[3]string{s + "inaH", s + "inoH", s + "InAm"},
				[3]string{s + "ini", s + "inoH", s + "izu"},
				[3]string{s + "i", s + "inI", s + "Ini"},
			)
		}
	} else if strings.HasSuffix(base, "u") {
		s := base[:len(base)-1]
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s + "uH", s + "U", s + "avaH"},
				[3]string{s + "um", s + "U", s + "Un"},
				[3]string{s + "unA", s + "uByAm", s + "uBiH"},
				[3]string{s + "ave", s + "uByAm", s + "uByaH"},
				[3]string{s + "oH", s + "uByAm", s + "uByaH"},
				[3]string{s + "oH", s + "voH", s + "UnAm"},
				[3]string{s + "O", s + "voH", s + "uzu"},
				[3]string{s + "o", s + "U", s + "avaH"},
			)
		} else if gender == "feminine" {
			res = buildGrid(
				[3]string{s + "uH", s + "U", s + "avaH"},
				[3]string{s + "um", s + "U", s + "UH"},
				[3]string{s + "vA", s + "uByAm", s + "uBiH"},
				[3]string{s + "vE", s + "uByAm", s + "uByaH"},
				[3]string{s + "vAH", s + "uByAm", s + "uByaH"},
				[3]string{s + "vAH", s + "voH", s + "UnAm"},
				[3]string{s + "vAm", s + "voH", s + "uzu"},
				[3]string{s + "o", s + "U", s + "avaH"},
			)
		} else if gender == "neuter" {
			res = buildGrid(
				[3]string{s + "u", s + "unI", s + "Uni"},
				[3]string{s + "u", s + "unI", s + "Uni"},
				[3]string{s + "unA", s + "uByAm", s + "uBiH"},
				[3]string{s + "une", s + "uByAm", s + "uByaH"},
				[3]string{s + "unaH", s + "uByAm", s + "uByaH"},
				[3]string{s + "unaH", s + "unoH", s + "UnAm"},
				[3]string{s + "uni", s + "unoH", s + "uzu"},
				[3]string{s + "u", s + "unI", s + "Uni"},
			)
		}
	} else if strings.HasSuffix(base, "U") {
		s := base[:len(base)-1]
		if gender == "feminine" {
			res = buildGrid(
				[3]string{s + "UH", s + "vO", s + "vaH"},
				[3]string{s + "Um", s + "vO", s + "UH"},
				[3]string{s + "vA", s + "UByaM", s + "UBiH"},
				[3]string{s + "vE", s + "UByaM", s + "UByaH"},
				[3]string{s + "vAH", s + "UByaM", s + "UByaH"},
				[3]string{s + "vAH", s + "voH", s + "UnAm"},
				[3]string{s + "vAm", s + "voH", s + "Uzu"},
				[3]string{s + "u", s + "vO", s + "vaH"},
			)
		}
	} else if strings.HasSuffix(base, "f") {
		s := base[:len(base)-1]
		if gender == "masculine" || gender == "feminine" {
			res = buildGrid(
				[3]string{s + "tA", s + "tArO", s + "tAraH"},
				[3]string{s + "tAram", s + "tArO", s + "tFn"},
				[3]string{s + "trA", s + "tfByAm", s + "tfBiH"},
				[3]string{s + "tre", s + "tfByAm", s + "tfByaH"},
				[3]string{s + "tuH", s + "tfByAm", s + "tfByaH"},
				[3]string{s + "tuH", s + "troH", s + "tFnAm"},
				[3]string{s + "tari", s + "troH", s + "tfzu"},
				[3]string{s + "taH", s + "tArO", s + "tAraH"},
			)
		}
	} else if strings.HasSuffix(base, "at") && len(base) >= 2 {
		s2 := base[:len(base)-2]
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
		} else if gender == "neuter" {
			res = buildGrid(
				[3]string{s2 + "at", s2 + "atI", s2 + "anti"},
				[3]string{s2 + "at", s2 + "atI", s2 + "anti"},
				[3]string{s2 + "atA", s2 + "adByAm", s2 + "adBiH"},
				[3]string{s2 + "ate", s2 + "adByAm", s2 + "adByaH"},
				[3]string{s2 + "ataH", s2 + "adByAm", s2 + "adByaH"},
				[3]string{s2 + "ataH", s2 + "atoH", s2 + "atAm"},
				[3]string{s2 + "ati", s2 + "atoH", s2 + "atsu"},
				[3]string{s2 + "at", s2 + "atI", s2 + "anti"},
			)
		}
	} else if strings.HasSuffix(base, "an") && len(base) >= 2 {
		s2 := base[:len(base)-2]
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
	} else if strings.HasSuffix(base, "in") && len(base) >= 2 {
		s2 := base[:len(base)-2]
		if gender == "masculine" {
			res = buildGrid(
				[3]string{s2 + "I", s2 + "inO", s2 + "inaH"},
				[3]string{s2 + "inam", s2 + "inO", s2 + "inaH"},
				[3]string{s2 + "inA", s2 + "iByAm", s2 + "iBiH"},
				[3]string{s2 + "ine", s2 + "iByAm", s2 + "iByaH"},
				[3]string{s2 + "inaH", s2 + "iByAm", s2 + "iByaH"},
				[3]string{s2 + "inaH", s2 + "inoH", s2 + "inAm"},
				[3]string{s2 + "ini", s2 + "inoH", s2 + "izu"},
				[3]string{s2 + "in", s2 + "inO", s2 + "inaH"},
			)
		}
	} else if strings.HasSuffix(base, "as") && len(base) >= 2 {
		s2 := base[:len(base)-2]
		if gender == "neuter" {
			res = buildGrid(
				[3]string{s2 + "aH", s2 + "asI", s2 + "AMsi"},
				[3]string{s2 + "aH", s2 + "asI", s2 + "AMsi"},
				[3]string{s2 + "asA", s2 + "oByAm", s2 + "oBiH"},
				[3]string{s2 + "ase", s2 + "oByAm", s2 + "oByaH"},
				[3]string{s2 + "asaH", s2 + "oByAm", s2 + "oByaH"},
				[3]string{s2 + "asaH", s2 + "asoH", s2 + "asAm"},
				[3]string{s2 + "asi", s2 + "asoH", s2 + "asu"},
				[3]string{s2 + "aH", s2 + "asI", s2 + "AMsi"},
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