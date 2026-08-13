package engine

import (
	"regexp"
)

type GuessedStem struct {
	Stem   string `json:"stem"`
	Gender string `json:"gender"`
	Case   string `json:"case"`
	Vacana string `json:"vacana"`
}

type stemRule struct {
	pattern     *regexp.Regexp
	replacement string
	gender      string
	caseName    string
	vacana      string
}

func rule(pat, rep, gen, c, vac string) stemRule {
	return stemRule{
		pattern:     regexp.MustCompile(pat),
		replacement: rep,
		gender:      gen,
		caseName:    c,
		vacana:      vac,
	}
}

var stemmingRules []stemRule

func init() {
	stemmingRules = []stemRule{
		// KINSHIP 'f' STEMS
		rule(`pitarO$`, "pitf", "masc", "prathama/dvitiya", "dvi"),
		rule(`pitaraH$`, "pitf", "masc", "prathama", "bahu"),
		rule(`pitaram$`, "pitf", "masc", "dvitiya", "eka"),
		rule(`pitFn$`, "pitf", "masc", "dvitiya", "bahu"),
		rule(`pitrA$`, "pitf", "masc", "tritiya", "eka"),
		rule(`pitre$`, "pitf", "masc", "caturthi", "eka"),
		rule(`pituH$`, "pitf", "masc", "panchami/sasthi", "eka"),
		rule(`pitari$`, "pitf", "masc", "saptami", "eka"),

		rule(`mAtarO$`, "mAtf", "fem", "prathama/dvitiya", "dvi"),
		rule(`mAtaraH$`, "mAtf", "fem", "prathama", "bahu"),
		rule(`mAtaram$`, "mAtf", "fem", "dvitiya", "eka"),
		rule(`mAtFH$`, "mAtf", "fem", "dvitiya", "bahu"),
		rule(`mAtrA$`, "mAtf", "fem", "tritiya", "eka"),
		rule(`mAtre$`, "mAtf", "fem", "caturthi", "eka"),
		rule(`mAtuH$`, "mAtf", "fem", "panchami/sasthi", "eka"),
		rule(`mAtari$`, "mAtf", "fem", "saptami", "eka"),

		// MONOSYLLABIC ROOT NOUNS ('dhI', 'BU', 'SrI')
		rule(`dhiyO$`, "dhI", "fem", "prathama/dvitiya", "dvi"),
		rule(`dhiyaH$`, "dhI", "fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`dhiyam$`, "dhI", "fem", "dvitiya", "eka"),
		rule(`dhiyA$`, "dhI", "fem", "tritiya", "eka"),
		rule(`dhiye$`, "dhI", "fem", "caturthi", "eka"),
		rule(`dhiyi$`, "dhI", "fem", "saptami", "eka"),

		rule(`BuvO$`, "BU", "fem", "prathama/dvitiya", "dvi"),
		rule(`BuvaH$`, "BU", "fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`Buvam$`, "BU", "fem", "dvitiya", "eka"),
		rule(`BuvA$`, "BU", "fem", "tritiya", "eka"),
		rule(`Buve$`, "BU", "fem", "caturthi", "eka"),
		rule(`Buvi$`, "BU", "fem", "saptami", "eka"),

		rule(`SriyO$`, "SrI", "fem", "prathama/dvitiya", "dvi"),
		rule(`SriyaH$`, "SrI", "fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`Sriyam$`, "SrI", "fem", "dvitiya", "eka"),
		rule(`SriyA$`, "SrI", "fem", "tritiya", "eka"),
		rule(`Sriye$`, "SrI", "fem", "caturthi", "eka"),
		rule(`Sriyi$`, "SrI", "fem", "saptami", "eka"),

		// MASCULINE 'as' & 'vas' STEMS
		rule(`vidvAn$`, "vidvas", "masc", "prathama", "eka"),
		rule(`vidvAMsO$`, "vidvas", "masc", "prathama/dvitiya", "dvi"),
		rule(`vidvAMsaH$`, "vidvas", "masc", "prathama", "bahu"),
		rule(`vidvAMsam$`, "vidvas", "masc", "dvitiya", "eka"),
		rule(`viduzaH$`, "vidvas", "masc", "dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`viduzA$`, "vidvas", "masc", "tritiya", "eka"),
		rule(`viduze$`, "vidvas", "masc", "caturthi", "eka"),
		rule(`viduzi$`, "vidvas", "masc", "saptami", "eka"),

		rule(`candramAH$`, "candramas", "masc", "prathama", "eka"),
		rule(`candramasO$`, "candramas", "masc", "prathama/dvitiya", "dvi"),
		rule(`candramasaH$`, "candramas", "masc", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`candramasam$`, "candramas", "masc", "dvitiya", "eka"),
		rule(`candramasA$`, "candramas", "masc", "tritiya", "eka"),

		// PRONOMINAL ADJECTIVES ('sarva')
		rule(`sarvasmE$`, "sarva", "masc/neut", "caturthi", "eka"),
		rule(`sarvasmAt$`, "sarva", "masc/neut", "panchami", "eka"),
		rule(`sarvezAm$`, "sarva", "masc/neut", "sasthi", "bahu"),
		rule(`sarvasmin$`, "sarva", "masc/neut", "saptami", "eka"),

		// 'a' STEMS
		rule(`aH$`, "a", "masc", "prathama", "eka"),
		rule(`O$`, "a", "masc", "prathama/dvitiya", "dvi"),
		rule(`AH$`, "a", "masc", "prathama", "bahu"),
		rule(`am$`, "a", "masc/neut", "dvitiya/prathama", "eka"),
		rule(`An$`, "a", "masc", "dvitiya", "bahu"),
		rule(`e[nR]a$`, "a", "masc/neut", "tritiya", "eka"),
		rule(`AByAm$`, "a", "masc/neut", "tritiya/caturthi/panchami", "dvi"),
		rule(`EH$`, "a", "masc/neut", "tritiya", "bahu"),
		rule(`Aya$`, "a", "masc/neut", "caturthi", "eka"),
		rule(`eByaH$`, "a", "masc/neut", "caturthi/panchami", "bahu"),
		rule(`At$`, "a", "masc/neut", "panchami", "eka"),
		rule(`asya$`, "a", "masc/neut", "sasthi", "eka"),
		rule(`ayoH$`, "a", "masc/neut", "sasthi/saptami", "dvi"),
		rule(`A[nR]Am$`, "a", "masc/neut", "sasthi", "bahu"),
		rule(`e$`, "a", "masc/neut", "saptami", "eka"),
		rule(`ezu$`, "a", "masc/neut", "saptami", "bahu"),
		rule(`Ani$`, "a", "neut", "prathama/dvitiya", "bahu"),

		// 'A' STEMS (FEMININE)
		rule(`A$`, "A", "fem", "prathama", "eka"),
		rule(`Am$`, "A", "fem", "dvitiya", "eka"),
		rule(`ayA$`, "A", "fem", "tritiya", "eka"),
		rule(`AyE$`, "A", "fem", "caturthi", "eka"),
		rule(`AyAH$`, "A", "fem", "panchami/sasthi", "eka"),
		rule(`AyAm$`, "A", "fem", "saptami", "eka"),

		// 'i' / 'I' STEMS
		rule(`iH$`, "i", "masc/fem", "prathama", "eka"),
		rule(`im$`, "i", "masc/fem/neut", "dvitiya/prathama", "eka"),
		rule(`In$`, "i", "masc", "dvitiya", "bahu"),
		rule(`IH$`, "i", "fem", "prathama/dvitiya", "bahu"),
		rule(`i[nR]A$`, "i", "masc/neut", "tritiya", "eka"),
		rule(`yA$`, "I", "fem", "tritiya", "eka"),
		rule(`aye$`, "i", "masc/fem", "caturthi", "eka"),
		rule(`yE$`, "I", "fem", "caturthi", "eka"),
		rule(`eH$`, "i", "masc/fem", "panchami/sasthi", "eka"),
		rule(`yAH$`, "I", "fem", "panchami/sasthi", "eka"),
		rule(`yAm$`, "I", "fem", "saptami", "eka"),
		rule(`I[nR]Am$`, "i", "masc/fem/neut", "sasthi", "bahu"),
		rule(`izu$`, "i", "masc/fem/neut", "saptami", "bahu"),
		rule(`Izu$`, "I", "fem", "saptami", "bahu"),
		rule(`iByAm$`, "i", "masc/fem/neut", "tritiya/caturthi/panchami", "dvi"),
		rule(`IByAm$`, "I", "fem", "tritiya/caturthi/panchami", "dvi"),

		// 'u' / 'U' STEMS
		rule(`uH$`, "u", "masc/fem", "prathama", "eka"),
		rule(`um$`, "u", "masc/fem/neut", "dvitiya/prathama", "eka"),
		rule(`Un$`, "u", "masc", "dvitiya", "bahu"),
		rule(`UH$`, "U", "fem", "prathama/dvitiya", "bahu"),
		rule(`u[nR]A$`, "u", "masc/neut", "tritiya", "eka"),
		rule(`vA$`, "u", "fem", "tritiya", "eka"),
		rule(`ave$`, "u", "masc/fem", "caturthi", "eka"),
		rule(`vE$`, "U", "fem", "caturthi", "eka"),
		rule(`oH$`, "u", "masc/fem", "panchami/sasthi", "eka"),
		rule(`vAH$`, "U", "fem", "panchami/sasthi", "eka"),
		rule(`vAm$`, "U", "fem", "saptami", "eka"),

		// AGENT 'f' STEMS
		rule(`tA$`, "tf", "masc/fem", "prathama", "eka"),
		rule(`tArO$`, "tf", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`tAraH$`, "tf", "masc/fem", "prathama", "bahu"),
		rule(`tAram$`, "tf", "masc/fem", "dvitiya", "eka"),
		rule(`tFn$`, "tf", "masc", "dvitiya", "bahu"),
		rule(`trA$`, "tf", "masc/fem", "tritiya", "eka"),
		rule(`tre$`, "tf", "masc/fem", "caturthi", "eka"),
		rule(`tuH$`, "tf", "masc/fem", "panchami/sasthi", "eka"),

		// HALANTA: 'c', 'j', 'd', 't', 'z' STEMS
		rule(`kzu$`, "c", "masc/fem", "saptami", "bahu"),
		rule(`cO$`, "c", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`caH$`, "c", "masc/fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`cam$`, "c", "masc/fem", "dvitiya", "eka"),
		rule(`cA$`, "c", "masc/fem", "tritiya", "eka"),
		rule(`ce$`, "c", "masc/fem", "caturthi", "eka"),
		rule(`ci$`, "c", "masc/fem", "saptami", "eka"),
		rule(`gByAm$`, "c", "masc/fem", "tritiya/caturthi/panchami", "dvi"),

		rule(`jO$`, "j", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`jaH$`, "j", "masc/fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`jam$`, "j", "masc/fem", "dvitiya", "eka"),
		rule(`jA$`, "j", "masc/fem", "tritiya", "eka"),
		rule(`je$`, "j", "masc/fem", "caturthi", "eka"),
		rule(`ji$`, "j", "masc/fem", "saptami", "eka"),

		rule(`tsu$`, "d", "masc/fem", "saptami", "bahu"),
		rule(`dO$`, "d", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`daH$`, "d", "masc/fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`dam$`, "d", "masc/fem", "dvitiya", "eka"),
		rule(`dA$`, "d", "masc/fem", "tritiya", "eka"),
		rule(`de$`, "d", "masc/fem", "caturthi", "eka"),
		rule(`di$`, "d", "masc/fem", "saptami", "eka"),
		rule(`dByAm$`, "d", "masc/fem", "tritiya/caturthi/panchami", "dvi"),

		rule(`w$`, "z", "masc/fem", "prathama", "eka"),
		rule(`wsu$`, "z", "masc/fem", "saptami", "bahu"),
		rule(`qByAm$`, "z", "masc/fem", "tritiya/caturthi/panchami", "dvi"),
		rule(`zO$`, "z", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`zaH$`, "z", "masc/fem", "prathama/dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`zam$`, "z", "masc/fem", "dvitiya", "eka"),

		// HALANTA: 'is' & 'us' NEUTER STEMS
		rule(`izI$`, "is", "neut", "prathama/dvitiya", "dvi"),
		rule(`IMzi$`, "is", "neut", "prathama/dvitiya", "bahu"),
		rule(`izA$`, "is", "neut", "tritiya", "eka"),
		rule(`irByAm$`, "is", "neut", "tritiya/caturthi/panchami", "dvi"),

		rule(`uzI$`, "us", "neut", "prathama/dvitiya", "dvi"),
		rule(`UMzi$`, "us", "neut", "prathama/dvitiya", "bahu"),
		rule(`uzA$`, "us", "neut", "tritiya", "eka"),
		rule(`urByAm$`, "us", "neut", "tritiya/caturthi/panchami", "dvi"),

		// HALANTA: 'vat' & 'mat' STEMS
		rule(`vAn$`, "vat", "masc", "prathama", "eka"),
		rule(`vantO$`, "vat", "masc", "prathama/dvitiya", "dvi"),
		rule(`vantaH$`, "vat", "masc", "prathama", "bahu"),
		rule(`vantam$`, "vat", "masc", "dvitiya", "eka"),
		rule(`vataH$`, "vat", "masc/neut", "dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`vatA$`, "vat", "masc/neut", "tritiya", "eka"),
		rule(`vadByAm$`, "vat", "masc/neut", "tritiya/caturthi/panchami", "dvi"),

		rule(`mAn$`, "mat", "masc", "prathama", "eka"),
		rule(`mantO$`, "mat", "masc", "prathama/dvitiya", "dvi"),
		rule(`mantaH$`, "mat", "masc", "prathama", "bahu"),
		rule(`mantam$`, "mat", "masc", "dvitiya", "eka"),
		rule(`mataH$`, "mat", "masc/neut", "dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`matA$`, "mat", "masc/neut", "tritiya", "eka"),
		rule(`madByAm$`, "mat", "masc/neut", "tritiya/caturthi/panchami", "dvi"),

		// HALANTA: 'an', 'in', 'at', 'as'
		rule(`an$`, "at", "masc", "prathama/sambodhana", "eka"),
		rule(`antO$`, "at", "masc", "prathama/dvitiya", "dvi"),
		rule(`antaH$`, "at", "masc", "prathama", "bahu"),
		rule(`antam$`, "at", "masc", "dvitiya", "eka"),
		rule(`ataH$`, "at", "masc/neut", "dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`nA$`, "an", "masc/neut", "tritiya", "eka"),
		rule(`inO$`, "in", "masc", "prathama/dvitiya", "dvi"),
		rule(`inaH$`, "in", "masc", "prathama/dvitiya/panchami/sasthi", "bahu"),
		rule(`aH$`, "as", "neut", "prathama/dvitiya", "eka"),
		rule(`asI$`, "as", "neut", "prathama/dvitiya", "dvi"),
		rule(`AMsi$`, "as", "neut", "prathama/dvitiya", "bahu"),
	}
}

// GetStems applies stemming rules to an inflected Sanskrit word to guess its potential stems
func GetStems(word string) []GuessedStem {
	var guessed []GuessedStem
	seen := make(map[string]bool)

	for _, r := range stemmingRules {
		if r.pattern.MatchString(word) {
			replaced := r.pattern.ReplaceAllString(word, r.replacement)
			key := replaced + "|" + r.gender + "|" + r.caseName + "|" + r.vacana
			if !seen[key] {
				seen[key] = true
				guessed = append(guessed, GuessedStem{
					Stem:   replaced,
					Gender: r.gender,
					Case:   r.caseName,
					Vacana: r.vacana,
				})
			}
		}
	}
	return guessed
}
