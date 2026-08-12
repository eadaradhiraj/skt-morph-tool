package engine

import "database/sql"

type pronounEntry struct {
	base   string
	gender string
	cCase  string
	vacana string
	form   string
}

var pronounMap []pronounEntry

func init() {
	addP := func(b, g, c, eka, dvi, bahu string) {
		pronounMap = append(pronounMap, pronounEntry{b, g, c, "eka", eka})
		if dvi != "" {
			pronounMap = append(pronounMap, pronounEntry{b, g, c, "dvi", dvi})
		}
		if bahu != "" {
			pronounMap = append(pronounMap, pronounEntry{b, g, c, "bahu", bahu})
		}
	}

	// TAD (Masculine)
	addP("tad", "masculine", "prathama", "saH", "tO", "te")
	addP("tad", "masculine", "dvitiya", "tam", "tO", "tAn")
	addP("tad", "masculine", "tritiya", "tena", "tAByAm", "tEH")
	addP("tad", "masculine", "caturthi", "tasmE", "tAByAm", "teByaH")
	addP("tad", "masculine", "panchami", "tasmAt", "tAByAm", "teByaH")
	addP("tad", "masculine", "sasthi", "tasya", "tayoH", "tezAm")
	addP("tad", "masculine", "saptami", "tasmin", "tayoH", "tezu")

	// TAD (Feminine)
	addP("tad", "feminine", "prathama", "sA", "te", "tAH")
	addP("tad", "feminine", "dvitiya", "tAm", "te", "tAH")
	addP("tad", "feminine", "tritiya", "tayA", "tAByAm", "tABiH")
	addP("tad", "feminine", "caturthi", "tasyE", "tAByAm", "tAByaH")
	addP("tad", "feminine", "panchami", "tasyAH", "tAByAm", "tAByaH")
	addP("tad", "feminine", "sasthi", "tasyAH", "tayoH", "tAsAm")
	addP("tad", "feminine", "saptami", "tasyAm", "tayoH", "tAsu")

	// TAD (Neuter)
	addP("tad", "neuter", "prathama", "tat", "te", "tAni")
	addP("tad", "neuter", "dvitiya", "tat", "te", "tAni")

	// KIM (Masculine)
	addP("kim", "masculine", "prathama", "kaH", "kO", "ke")
	addP("kim", "masculine", "dvitiya", "kam", "kO", "kAn")
	addP("kim", "masculine", "caturthi", "kasmE", "kAByAm", "keByaH")

	// ASMAD (First Person)
	addP("asmad", "any", "prathama", "aham", "AvAm", "vayam")
	addP("asmad", "any", "dvitiya", "mAm", "AvAm", "asmAn")
	addP("asmad", "any", "tritiya", "mayA", "AvAByAm", "asmAByiH")
	addP("asmad", "any", "caturthi", "mahyam", "AvAByAm", "asmAByam")
	addP("asmad", "any", "panchami", "mat", "AvAByAm", "asmat")
	addP("asmad", "any", "sasthi", "mama", "AvayoH", "asmAkam")
	addP("asmad", "any", "saptami", "mayi", "AvayoH", "asmAsu")

	// YUZMAD (Second Person)
	addP("yuzmad", "any", "prathama", "tvam", "yuvAm", "yUyam")
	addP("yuzmad", "any", "dvitiya", "tvAm", "yuvAm", "yuzmAn")
	addP("yuzmad", "any", "tritiya", "tvayA", "yuvAByAm", "yuzmABiH")
	addP("yuzmad", "any", "caturthi", "tubhyam", "yuvAByAm", "yuzmAByam")
	addP("yuzmad", "any", "panchami", "tvat", "yuvAByAm", "yuzmat")
	addP("yuzmad", "any", "sasthi", "tava", "yuvayoH", "yuzmAkam")
	addP("yuzmad", "any", "saptami", "tvayi", "yuvayoH", "yuzmAsu")
}

func AnalyzePronoun(db *sql.DB, word string) []map[string]any {
	var results []map[string]any

	for _, entry := range pronounMap {
		if entry.form == word {
			results = append(results, map[string]any{
				"type":      "pronoun",
				"base_form": entry.base,
				"gender":    entry.gender,
				"case":      entry.cCase,
				"vacana":    entry.vacana,
			})
		}
	}

	return results
}