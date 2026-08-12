package engine

import (
	"database/sql"
	"fmt"
	"strings"
)

func exactMatch(rowVal, word string) bool {
	parts := strings.Split(rowVal, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == word {
			return true
		}
	}
	return false
}

func isValidParticiple(db *sql.DB, dhatuID, upasarga, pratyaya string) bool {
	if pratyaya == "lyap" && upasarga == "" {
		return false
	}
	atmanepadaPratyayas := []string{"SAnac", "cAnaS", "sya-SAnac", "BAvakarma-SAnac", "sya-BAvakarma-SAnac"}
	isAtmanepada := false
	for _, ap := range atmanepadaPratyayas {
		if pratyaya == ap {
			isAtmanepada = true
			break
		}
	}

	if isAtmanepada {
		var dummy int
		query := "SELECT 1 FROM conjugations WHERE dhatu_id = ?1 AND upasarga = ?2 AND voice = 'Atmanepadam' LIMIT 1"
		err := db.QueryRow(query, dhatuID, upasarga).Scan(&dummy)
		if err != nil {
			return false
		}
	}
	return true
}

func fetchVerbs(db *sql.DB, word string) []map[string]any {
	var results []map[string]any
	searchTerm := fmt.Sprintf("%%%s%%", word)

	query := `SELECT dhatu_id, upasarga, derivative, prayoga, lakara, voice, purusha, eka, dvi, bahu 
	          FROM conjugations WHERE eka LIKE ?1 OR dvi LIKE ?1 OR bahu LIKE ?1`

	rows, err := db.Query(query, searchTerm)
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var dhatuID, upasarga, derivative, prayoga, lakara, voice, purusha, eka, dvi, bahu sql.NullString
		if err := rows.Scan(&dhatuID, &upasarga, &derivative, &prayoga, &lakara, &voice, &purusha, &eka, &dvi, &bahu); err != nil {
			continue
		}

		var matchedVacanas []string
		if exactMatch(eka.String, word) {
			matchedVacanas = append(matchedVacanas, "eka")
		}
		if exactMatch(dvi.String, word) {
			matchedVacanas = append(matchedVacanas, "dvi")
		}
		if exactMatch(bahu.String, word) {
			matchedVacanas = append(matchedVacanas, "bahu")
		}

		for _, vacana := range matchedVacanas {
			results = append(results, map[string]any{
				"type":       "verb",
				"dhatu_id":   dhatuID.String,
				"upasarga":   upasarga.String,
				"derivative": derivative.String,
				"prayoga":    prayoga.String,
				"lakara":     lakara.String,
				"voice":      voice.String,
				"purusha":    purusha.String,
				"vacana":     vacana,
			})
		}
	}
	return results
}

func analyzeVerb(db *sql.DB, word string) []map[string]any {
	results := fetchVerbs(db, word)
	if len(results) == 0 {
		for _, split := range GetUpasargaSplits(word) {
			upa, stripped := split[0], split[1]
			subResults := fetchVerbs(db, stripped)
			for _, res := range subResults {
				if res["upasarga"] == "" {
					res["upasarga"] = upa
					res["note"] = "Dynamically matched via Sandhi split"
					results = append(results, res)
				}
			}
			if len(results) > 0 {
				break
			}
		}
	}
	return results
}

func fetchParticiples(db *sql.DB, word string) []map[string]any {
	var results []map[string]any
	searchTerm := fmt.Sprintf("%%%s%%", word)

	query := `SELECT dhatu_id, upasarga, derivative, pratyaya, base_form, masculine, feminine, neuter 
	          FROM participles WHERE base_form LIKE ?1 OR masculine LIKE ?1 OR feminine LIKE ?1 OR neuter LIKE ?1`

	rows, err := db.Query(query, searchTerm)
	if err != nil {
		return results
	}
	defer rows.Close()

	avyayaPratyayas := []string{"tumun", "ktvA", "lyap", "Ramul"}

	for rows.Next() {
		var dhatuID, upasarga, derivative, pratyaya, baseForm, masc, fem, neut sql.NullString
		if err := rows.Scan(&dhatuID, &upasarga, &derivative, &pratyaya, &baseForm, &masc, &fem, &neut); err != nil {
			continue
		}

		if !isValidParticiple(db, dhatuID.String, upasarga.String, pratyaya.String) {
			continue
		}

		var matchedCols []string
		if exactMatch(baseForm.String, word) {
			matchedCols = append(matchedCols, "base_form")
		}
		if exactMatch(masc.String, word) {
			matchedCols = append(matchedCols, "masculine")
		}
		if exactMatch(fem.String, word) {
			matchedCols = append(matchedCols, "feminine")
		}
		if exactMatch(neut.String, word) {
			matchedCols = append(matchedCols, "neuter")
		}

		isAvyaya := false
		for _, av := range avyayaPratyayas {
			if pratyaya.String == av {
				isAvyaya = true
				break
			}
		}

		pType := "participle"
		if isAvyaya {
			pType = "avyaya"
		}

		for _, col := range matchedCols {
			pJSON := map[string]any{
				"type":       pType,
				"dhatu_id":   dhatuID.String,
				"upasarga":   upasarga.String,
				"derivative": derivative.String,
				"pratyaya":   pratyaya.String,
				"base_form":  baseForm.String,
			}

			if col == "masculine" || col == "feminine" || col == "neuter" {
				pJSON["gender"] = col
				pJSON["case"] = "prathama"
				pJSON["vacana"] = "eka"
			} else {
				pJSON["note"] = "Matched uninflected base form"
			}
			results = append(results, pJSON)
		}
	}
	return results
}

func analyzeParticiple(db *sql.DB, word string) []map[string]any {
	results := fetchParticiples(db, word)
	if len(results) == 0 {
		for _, split := range GetUpasargaSplits(word) {
			upa, stripped := split[0], split[1]
			subResults := fetchParticiples(db, stripped)
			for _, res := range subResults {
				if res["upasarga"] == "" {
					res["upasarga"] = upa
					res["note"] = "Dynamically matched via Sandhi split"
					results = append(results, res)
				}
			}
			if len(results) > 0 {
				break
			}
		}
	}
	return results
}

func analyzeDeclension(db *sql.DB, word string) []map[string]any {
	guessedStems := GetStems(word)
	var results []map[string]any

	for _, guess := range guessedStems {
		stem := guess.Stem
		searchTerm := fmt.Sprintf("%%%s%%", stem)

		query := "SELECT dhatu_id, upasarga, pratyaya, base_form FROM participles WHERE base_form LIKE ?1"
		rows, err := db.Query(query, searchTerm)

		var exactMatches []map[string]string
		if err == nil {
			for rows.Next() {
				var dhatuID, upasarga, pratyaya, baseForm sql.NullString
				if err := rows.Scan(&dhatuID, &upasarga, &pratyaya, &baseForm); err == nil {
					if exactMatch(baseForm.String, stem) {
						exactMatches = append(exactMatches, map[string]string{
							"dhatu_id": dhatuID.String,
							"upasarga": upasarga.String,
							"pratyaya": pratyaya.String,
						})
					}
				}
			}
			rows.Close()
		}

		if len(exactMatches) > 0 {
			for _, m := range exactMatches {
				r := map[string]any{
					"type":      "declension",
					"stem":      guess.Stem,
					"gender":    guess.Gender,
					"case":      guess.Case,
					"vacana":    guess.Vacana,
					"base_form": stem,
					"dhatu_id":  m["dhatu_id"],
					"upasarga":  m["upasarga"],
					"pratyaya":  m["pratyaya"],
				}
				results = append(results, r)
			}
		} else {
			foundDynamic := false
			for _, split := range GetUpasargaSplits(stem) {
				upa, strippedStem := split[0], split[1]
				sTerm := fmt.Sprintf("%%%s%%", strippedStem)

				dQuery := "SELECT dhatu_id, pratyaya, base_form FROM participles WHERE base_form LIKE ?1 AND upasarga = ''"
				dRows, err := db.Query(dQuery, sTerm)
				if err == nil {
					for dRows.Next() {
						var dhatuID, pratyaya, baseForm sql.NullString
						if err := dRows.Scan(&dhatuID, &pratyaya, &baseForm); err == nil {
							if exactMatch(baseForm.String, strippedStem) {
								r := map[string]any{
									"type":      "declension",
									"stem":      guess.Stem,
									"gender":    guess.Gender,
									"case":      guess.Case,
									"vacana":    guess.Vacana,
									"base_form": stem,
									"dhatu_id":  dhatuID.String,
									"upasarga":  upa,
									"pratyaya":  pratyaya.String,
									"note":      "Dynamic Upasarga Match",
								}
								results = append(results, r)
								foundDynamic = true
							}
						}
					}
					dRows.Close()
				}
				if foundDynamic {
					break
				}
			}

			if !foundDynamic {
				r := map[string]any{
					"type":      "declension",
					"stem":      guess.Stem,
					"gender":    guess.Gender,
					"case":      guess.Case,
					"vacana":    guess.Vacana,
					"base_form": stem,
					"dhatu_id":  nil,
					"upasarga":  nil,
					"pratyaya":  nil,
				}
				results = append(results, r)
			}
		}
	}
	return results
}

// Analyze is the master orchestrator for morphological analysis of a Sanskrit word
func Analyze(db *sql.DB, word string) map[string]any {
	verbs := analyzeVerb(db, word)
	participles := analyzeParticiple(db, word)
	declensions := analyzeDeclension(db, word)
	pronouns := AnalyzePronoun(db, word)
	numerals := AnalyzeNumeral(db, word)
	irregulars := AnalyzeIrregular(db, word)

	return map[string]any{
		"searched_word": word,
		"verbs":         verbs,
		"participles":   participles,
		"declensions":   declensions,
		"pronouns":      pronouns,
		"numerals":      numerals,
		"irregulars":    irregulars,
	}
}