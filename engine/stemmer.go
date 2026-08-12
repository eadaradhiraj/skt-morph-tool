package engine

import (
	"regexp"
)

// GuessedStem represents a stem guessed from an inflected word ending
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
	caseName    string // 'case' is a reserved keyword in Go
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

// init() runs automatically when the package is imported (replaces Rust's lazy_static!)
func init() {
	stemmingRules = []stemRule{
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
		rule(`iH$`, "i", "masc/fem", "prathama", "eka"),
		rule(`im$`, "i", "masc/fem/neut", "dvitiya/prathama", "eka"),
		rule(`In$`, "i", "masc", "dvitiya", "bahu"),
		rule(`IH$`, "i", "fem", "prathama/dvitiya", "bahu"),
		rule(`i[nR]A$`, "i", "masc/neut", "tritiya", "eka"),
		rule(`yA$`, "i", "fem", "tritiya", "eka"),
		rule(`aye$`, "i", "masc/fem", "caturthi", "eka"),
		rule(`yE$`, "i", "fem", "caturthi", "eka"),
		rule(`eH$`, "i", "masc/fem", "panchami/sasthi", "eka"),
		rule(`yAH$`, "i", "fem", "panchami/sasthi", "eka"),
		rule(`yAm$`, "i", "fem", "saptami", "eka"),
		rule(`I[nR]Am$`, "i", "masc/fem/neut", "sasthi", "bahu"),
		rule(`izu$`, "i", "masc/fem/neut", "saptami", "bahu"),
		rule(`iByAm$`, "i", "masc/fem/neut", "tritiya/caturthi/panchami", "dvi"),
		rule(`iByaH$`, "i", "masc/fem/neut", "caturthi/panchami", "bahu"),
		rule(`iBiH$`, "i", "masc/fem/neut", "tritiya", "bahu"),
		rule(`I$`, "i", "masc/fem/neut", "prathama/dvitiya", "dvi"),
		rule(`ayaH$`, "i", "masc/fem", "prathama", "bahu"),
		rule(`Ini$`, "i", "neut", "prathama/dvitiya", "bahu"),
		rule(`uH$`, "u", "masc/fem", "prathama", "eka"),
		rule(`um$`, "u", "masc/fem/neut", "dvitiya/prathama", "eka"),
		rule(`Un$`, "u", "masc", "dvitiya", "bahu"),
		rule(`UH$`, "u", "fem", "prathama/dvitiya", "bahu"),
		rule(`u[nR]A$`, "u", "masc/neut", "tritiya", "eka"),
		rule(`vA$`, "u", "fem", "tritiya", "eka"),
		rule(`ave$`, "u", "masc/fem", "caturthi", "eka"),
		rule(`vE$`, "u", "fem", "caturthi", "eka"),
		rule(`oH$`, "u", "masc/fem", "panchami/sasthi", "eka"),
		rule(`vAH$`, "u", "fem", "panchami/sasthi", "eka"),
		rule(`vAm$`, "u", "fem", "saptami", "eka"),
		rule(`U[nR]Am$`, "u", "masc/fem/neut", "sasthi", "bahu"),
		rule(`uzu$`, "u", "masc/fem/neut", "saptami", "bahu"),
		rule(`uByAm$`, "u", "masc/fem/neut", "tritiya/caturthi/panchami", "dvi"),
		rule(`uByaH$`, "u", "masc/fem/neut", "caturthi/panchami", "bahu"),
		rule(`uBiH$`, "u", "masc/fem/neut", "tritiya", "bahu"),
		rule(`U$`, "u", "masc/fem/neut", "prathama/dvitiya", "dvi"),
		rule(`avaH$`, "u", "masc/fem", "prathama", "bahu"),
		rule(`Uni$`, "u", "neut", "prathama/dvitiya", "bahu"),
		rule(`tA$`, "tf", "masc/fem", "prathama", "eka"),
		rule(`tArO$`, "tf", "masc/fem", "prathama/dvitiya", "dvi"),
		rule(`tAraH$`, "tf", "masc/fem", "prathama", "bahu"),
		rule(`tAram$`, "tf", "masc/fem", "dvitiya", "eka"),
		rule(`tFn$`, "tf", "masc", "dvitiya", "bahu"),
		rule(`tFH$`, "tf", "fem", "dvitiya", "bahu"),
		rule(`trA$`, "tf", "masc/fem", "tritiya", "eka"),
		rule(`tre$`, "tf", "masc/fem", "caturthi", "eka"),
		rule(`tuH$`, "tf", "masc/fem", "panchami/sasthi", "eka"),
		rule(`tari$`, "tf", "masc/fem", "saptami", "eka"),
		rule(`tF[nR]Am$`, "tf", "masc/fem", "sasthi", "bahu"),
		rule(`tfzu$`, "tf", "masc/fem", "saptami", "bahu"),
		rule(`tfByAm$`, "tf", "masc/fem", "tritiya/caturthi/panchami", "dvi"),
		rule(`tfByaH$`, "tf", "masc/fem", "caturthi/panchami", "bahu"),
		rule(`tfBiH$`, "tf", "masc/fem", "tritiya", "bahu"),

		// HALANTA: 'at' STEMS
		rule(`an$`, "at", "masc", "prathama/sambodhana", "eka"),
		rule(`antO$`, "at", "masc", "prathama/dvitiya", "dvi"),
		rule(`antaH$`, "at", "masc", "prathama", "bahu"),
		rule(`antam$`, "at", "masc", "dvitiya", "eka"),
		rule(`ataH$`, "at", "masc/neut", "dvitiya/panchami/sasthi", "bahu/eka"),
		rule(`atA$`, "at", "masc/neut", "tritiya", "eka"),
		rule(`ate$`, "at", "masc/neut", "caturthi", "eka"),
		rule(`ati$`, "at", "masc/neut", "saptami", "eka"),
		rule(`atoH$`, "at", "masc/neut", "sasthi/saptami", "dvi"),
		rule(`atAm$`, "at", "masc/neut", "sasthi", "bahu"),
		rule(`atsu$`, "at", "masc/neut", "saptami", "bahu"),
		rule(`adByAm$`, "at", "masc/neut", "tritiya/caturthi/panchami", "dvi"),
		rule(`adBiH$`, "at", "masc/neut", "tritiya", "bahu"),
		rule(`adByaH$`, "at", "masc/neut", "caturthi/panchami", "bahu"),
		rule(`at$`, "at", "neut", "prathama/dvitiya", "eka"),
		rule(`atI$`, "at", "neut", "prathama/dvitiya", "dvi"),
		rule(`anti$`, "at", "neut", "prathama/dvitiya", "bahu"),

		// CONSONANTS (Halanta)
		rule(`A$`, "an", "masc", "prathama", "eka"),
		rule(`AnO$`, "an", "masc", "prathama/dvitiya", "dvi"),
		rule(`AnaH$`, "an", "masc", "prathama", "bahu"),
		rule(`Anam$`, "an", "masc", "dvitiya", "eka"),
		rule(`nA$`, "an", "masc/neut", "tritiya", "eka"),
		rule(`I$`, "in", "masc", "prathama", "eka"),
		rule(`inO$`, "in", "masc", "prathama/dvitiya", "dvi"),
		rule(`inaH$`, "in", "masc", "prathama/dvitiya/panchami/sasthi", "bahu"),
		rule(`aH$`, "as", "neut", "prathama/dvitiya", "eka"),
		rule(`asI$`, "as", "neut", "prathama/dvitiya", "dvi"),
		rule(`AMsi$`, "as", "neut", "prathama/dvitiya", "bahu"),
		rule(`asA$`, "as", "masc/neut", "tritiya", "eka"),
		rule(`iH$`, "is", "neut", "prathama/dvitiya", "eka"),
		rule(`izI$`, "is", "neut", "prathama/dvitiya", "dvi"),
		rule(`IMzi$`, "is", "neut", "prathama/dvitiya", "bahu"),
		rule(`izA$`, "is", "masc/neut", "tritiya", "eka"),
		rule(`irBiH$`, "is", "masc/neut", "tritiya", "bahu"),
		rule(`uH$`, "us", "neut", "prathama/dvitiya", "eka"),
		rule(`uzI$`, "us", "neut", "prathama/dvitiya", "dvi"),
		rule(`UMzi$`, "us", "neut", "prathama/dvitiya", "bahu"),
		rule(`uzA$`, "us", "masc/neut", "tritiya", "eka"),
		rule(`urBiH$`, "us", "masc/neut", "tritiya", "bahu"),
		rule(`vAn$`, "vat", "masc", "prathama", "eka"),
		rule(`vantO$`, "vat", "masc", "prathama/dvitiya", "dvi"),
		rule(`vantaH$`, "vat", "masc", "prathama", "bahu"),
		rule(`vatA$`, "vat", "masc/neut", "tritiya", "eka"),
		rule(`van$`, "vat", "masc", "sambodhana", "eka"),
		rule(`mAn$`, "mat", "masc", "prathama", "eka"),
		rule(`mantO$`, "mat", "masc", "prathama/dvitiya", "dvi"),
		rule(`mantaH$`, "mat", "masc", "prathama", "bahu"),
		rule(`matA$`, "mat", "masc/neut", "tritiya", "eka"),
		rule(`man$`, "mat", "masc", "sambodhana", "eka"),
	}
}

// GetStems applies stemming rules to an inflected Sanskrit word to guess its potential stems
func GetStems(word string) []GuessedStem {
	var guessed []GuessedStem
	for _, r := range stemmingRules {
		if r.pattern.MatchString(word) {
			replaced := r.pattern.ReplaceAllString(word, r.replacement)
			guessed = append(guessed, GuessedStem{
				Stem:   replaced,
				Gender: r.gender,
				Case:   r.caseName,
				Vacana: r.vacana,
			})
		}
	}
	return guessed
}