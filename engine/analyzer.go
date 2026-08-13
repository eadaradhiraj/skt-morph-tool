package engine

import (
	"database/sql"
	"strings"
)

func exactMatch(rowVal, word string) bool {
	if rowVal == "" || word == "" {
		return false
	}
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

func fetchUpasargaMeaning(db *sql.DB, dhatuID, upasarga string) string {
	if upasarga == "" || dhatuID == "" {
		return ""
	}
	var meaning sql.NullString
	q := "SELECT meaning FROM upasarga_meanings WHERE dhatu_id = ?1 AND (upasarga_combination = ?2 OR upasarga_combination = REPLACE(?2, ' + ', '+') OR upasarga_combination = REPLACE(?2, '+', ' + ')) LIMIT 1"
	err := db.QueryRow(q, dhatuID, upasarga).Scan(&meaning)
	if err == nil && meaning.Valid && meaning.String != "" {
		return meaning.String
	}

	q2 := "SELECT meaning FROM upasargachandrika WHERE dhatu_id = ?1 AND (upasarga_combination = ?2 OR upasarga_combination = REPLACE(?2, ' + ', '+') OR upasarga_combination = REPLACE(?2, '+', ' + ')) LIMIT 1"
	err2 := db.QueryRow(q2, dhatuID, upasarga).Scan(&meaning)
	if err2 == nil && meaning.Valid && meaning.String != "" {
		return meaning.String
	}
	return ""
}

func fetchLiteraryAttestation(db *sql.DB, dhatuID, word string) map[string]string {
	if word == "" || dhatuID == "" {
		return nil
	}
	var book, text sql.NullString
	q := "SELECT book, literature_text FROM prayoga WHERE dhatu_id = ?1 AND form = ?2 LIMIT 1"
	err := db.QueryRow(q, dhatuID, word).Scan(&book, &text)
	if err == nil && (book.Valid || text.Valid) {
		return map[string]string{
			"book": book.String,
			"text": text.String,
		}
	}
	return nil
}

func fetchVerbs(db *sql.DB, word string) []map[string]any {
	var results []map[string]any

	query := `SELECT c.dhatu_id, c.upasarga, c.derivative, c.prayoga, c.lakara, c.voice, c.purusha, c.eka, c.dvi, c.bahu,
	                 COALESCE(i1.value, '') as root,
	                 COALESCE(i2.value, '') as meaning
	          FROM conjugations c
	          LEFT JOIN info i1 ON i1.dhatu_id = c.dhatu_id AND i1.key_name IN ('mUlaDAtuH', 'DAtuH')
	          LEFT JOIN info i2 ON i2.dhatu_id = c.dhatu_id AND i2.key_name IN ('aTfaH', 'meaning', 'eng', 'hin')
	          WHERE c.eka = ?1 OR c.eka LIKE ?1 || ',%' OR c.eka LIKE '%,' || ?1 OR c.eka LIKE '%,' || ?1 || ',%'
	             OR c.dvi = ?1 OR c.dvi LIKE ?1 || ',%' OR c.dvi LIKE '%,' || ?1 OR c.dvi LIKE '%,' || ?1 || ',%'
	             OR c.bahu = ?1 OR c.bahu LIKE ?1 || ',%' OR c.bahu LIKE '%,' || ?1 OR c.bahu LIKE '%,' || ?1 || ',%'`

	rows, err := db.Query(query, word)
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var dhatuID, upasarga, derivative, prayoga, lakara, voice, purusha, eka, dvi, bahu, root, meaning sql.NullString
		if err := rows.Scan(&dhatuID, &upasarga, &derivative, &prayoga, &lakara, &voice, &purusha, &eka, &dvi, &bahu, &root, &meaning); err != nil {
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

		if len(matchedVacanas) == 0 {
			continue
		}

		prefixedMeaning := fetchUpasargaMeaning(db, dhatuID.String, upasarga.String)
		literaryExample := fetchLiteraryAttestation(db, dhatuID.String, word)

		for _, vacana := range matchedVacanas {
			res := map[string]any{
				"type":       "verb",
				"dhatu_id":   dhatuID.String,
				"root":       root.String,
				"meaning":    meaning.String,
				"upasarga":   upasarga.String,
				"derivative": derivative.String,
				"prayoga":    prayoga.String,
				"lakara":     lakara.String,
				"voice":      voice.String,
				"purusha":    purusha.String,
				"vacana":     vacana,
			}
			if prefixedMeaning != "" {
				res["prefixed_meaning"] = prefixedMeaning
			}
			if literaryExample != nil {
				res["literature_attestation"] = literaryExample
			}
			results = append(results, res)
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
				if res["upasarga"] == "" || res["upasarga"] == nil {
					res["upasarga"] = upa
					res["note"] = "Dynamically matched via Sandhi split"
					if dhatuID, ok := res["dhatu_id"].(string); ok {
						pm := fetchUpasargaMeaning(db, dhatuID, upa)
						if pm != "" {
							res["prefixed_meaning"] = pm
						}
					}
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

	query := `SELECT dhatu_id, upasarga, derivative, pratyaya, base_form, masculine, feminine, neuter 
	          FROM participles 
	          WHERE base_form = ?1 OR base_form LIKE ?1 || ',%' OR base_form LIKE '%,' || ?1 OR base_form LIKE '%,' || ?1 || ',%'
	             OR masculine = ?1 OR masculine LIKE ?1 || ',%' OR masculine LIKE '%,' || ?1 OR masculine LIKE '%,' || ?1 || ',%'
	             OR feminine = ?1 OR feminine LIKE ?1 || ',%' OR feminine LIKE '%,' || ?1 OR feminine LIKE '%,' || ?1 || ',%'
	             OR neuter = ?1 OR neuter LIKE ?1 || ',%' OR neuter LIKE '%,' || ?1 OR neuter LIKE '%,' || ?1 || ',%'`

	rows, err := db.Query(query, word)
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

		if len(matchedCols) == 0 {
			continue
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

		prefixedMeaning := fetchUpasargaMeaning(db, dhatuID.String, upasarga.String)

		for _, col := range matchedCols {
			pJSON := map[string]any{
				"type":       pType,
				"dhatu_id":   dhatuID.String,
				"upasarga":   upasarga.String,
				"derivative": derivative.String,
				"pratyaya":   pratyaya.String,
				"base_form":  baseForm.String,
			}

			if prefixedMeaning != "" {
				pJSON["prefixed_meaning"] = prefixedMeaning
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
				if res["upasarga"] == "" || res["upasarga"] == nil {
					res["upasarga"] = upa
					res["note"] = "Dynamically matched via Sandhi split"
					if dhatuID, ok := res["dhatu_id"].(string); ok {
						pm := fetchUpasargaMeaning(db, dhatuID, upa)
						if pm != "" {
							res["prefixed_meaning"] = pm
						}
					}
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
	seen := make(map[string]bool)

	for _, guess := range guessedStems {
		stem := guess.Stem

		query := "SELECT dhatu_id, upasarga, pratyaya, base_form FROM participles WHERE base_form = ?1"
		rows, err := db.Query(query, stem)

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
				key := stem + "|" + guess.Gender + "|" + guess.Case + "|" + guess.Vacana + "|" + m["dhatu_id"]
				if seen[key] {
					continue
				}
				seen[key] = true

				pm := fetchUpasargaMeaning(db, m["dhatu_id"], m["upasarga"])
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
				if pm != "" {
					r["prefixed_meaning"] = pm
				}
				results = append(results, r)
			}
		} else {
			foundDynamic := false
			for _, split := range GetUpasargaSplits(stem) {
				upa, strippedStem := split[0], split[1]

				dQuery := "SELECT dhatu_id, pratyaya, base_form FROM participles WHERE base_form = ?1 AND upasarga = ''"
				dRows, err := db.Query(dQuery, strippedStem)
				if err == nil {
					for dRows.Next() {
						var dhatuID, pratyaya, baseForm sql.NullString
						if err := dRows.Scan(&dhatuID, &pratyaya, &baseForm); err == nil {
							if exactMatch(baseForm.String, strippedStem) {
								key := stem + "|" + guess.Gender + "|" + guess.Case + "|" + guess.Vacana + "|" + dhatuID.String
								if seen[key] {
									continue
								}
								seen[key] = true

								pm := fetchUpasargaMeaning(db, dhatuID.String, upa)
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
								if pm != "" {
									r["prefixed_meaning"] = pm
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
		}
	}
	return results
}

func Analyze(db *sql.DB, word string) map[string]any {
	word = strings.TrimSpace(word)
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
