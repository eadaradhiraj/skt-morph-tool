package engine

import "database/sql"

type pronounEntry struct {
	base   string
	gender string
	cCase  string
	forms  [3]string
}

var pronounMap = []pronounEntry{
	{"tad", "masculine", "prathama", [3]string{"saH", "tO", "te"}},
	{"tad", "masculine", "dvitiya", [3]string{"tam", "tO", "tAn"}},
	{"tad", "masculine", "caturthi", [3]string{"tasmE", "tAByAm", "teByaH"}},
	{"kim", "masculine", "prathama", [3]string{"kaH", "kO", "ke"}},
	{"asmad", "any", "prathama", [3]string{"aham", "AvAm", "vayam"}},
	{"yuzmad", "any", "prathama", [3]string{"tvam", "yuvAm", "yUyam"}},
}

// AnalyzePronoun analyzes pronouns against the hardcoded pronoun map
func AnalyzePronoun(db *sql.DB, word string) []map[string]any {
	var results []map[string]any
	vacanaList := []string{"eka", "dvi", "bahu"}

	for _, entry := range pronounMap {
		for i, form := range entry.forms {
			if form == word {
				results = append(results, map[string]any{
					"type":      "pronoun",
					"base_form": entry.base,
					"gender":    entry.gender,
					"case":      entry.cCase,
					"vacana":    vacanaList[i],
				})
			}
		}
	}

	return results
}