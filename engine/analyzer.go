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

func cleanRoot(root string) string {
	if root == "" {
		return root
	}
	// strip anubandha markers like ~ and trailing capital anubandha (e.g. dfN -> df, BU stays BU)
	root = strings.TrimSuffix(root, "~")
	// common anubandha: trailing N, M, Y etc. Strip if root ends with single uppercase and length>2
	// dfN -> df (N is anubandha), but keep BU (2 chars) as is
	if len(root) > 2 && root[len(root)-1] >= 'A' && root[len(root)-1] <= 'Z' {
		// check if previous char is not uppercase (avoid stripping BU where both are upper?)
		prev := root[len(root)-2]
		if !(prev >= 'A' && prev <= 'Z') {
			return root[:len(root)-1]
		}
	}
	return root
}

func hasTable(db *sql.DB, name string) bool {
	var n string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?1", name).Scan(&n)
	return err == nil && n == name
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
		if hasTable(db, "conjugations") {
			query := "SELECT 1 FROM conjugations WHERE dhatu_id = ?1 AND upasarga = ?2 AND voice = 'Atmanepadam' LIMIT 1"
			err := db.QueryRow(query, dhatuID, upasarga).Scan(&dummy)
			if err != nil {
				return false
			}
		} else if hasTable(db, "conjugation_forms") {
			// new schema: check if any Atmanepada form exists (form_type starting with 'a')
			query := "SELECT 1 FROM conjugation_forms WHERE dhatu_id = ?1 AND prefix = ?2 AND form_type LIKE 'a%' LIMIT 1"
			err := db.QueryRow(query, dhatuID, upasarga).Scan(&dummy)
			if err != nil {
				return false
			}
		}
	}
	return true
}

func fetchUpasargaMeaning(db *sql.DB, dhatuID, upasarga string) string {
	if upasarga == "" || dhatuID == "" {
		return ""
	}
	var meaning sql.NullString
	// old schema
	if hasTable(db, "upasarga_meanings") {
		q := "SELECT meaning FROM upasarga_meanings WHERE dhatu_id = ?1 AND (upasarga_combination = ?2 OR upasarga_combination = REPLACE(?2, ' + ', '+') OR upasarga_combination = REPLACE(?2, '+', ' + ')) LIMIT 1"
		err := db.QueryRow(q, dhatuID, upasarga).Scan(&meaning)
		if err == nil && meaning.Valid && meaning.String != "" {
			return meaning.String
		}
	} else if hasTable(db, "upasarga_entries") {
		// new schema: upasarga_entries / upasarga_prefix_meta store artha_hindi
		q := "SELECT artha_hindi FROM upasarga_entries WHERE dhatu_id = ?1 AND prefix = ?2 LIMIT 1"
		err := db.QueryRow(q, dhatuID, upasarga).Scan(&meaning)
		if err == nil && meaning.Valid && meaning.String != "" {
			return meaning.String
		}
		q2 := "SELECT artha_hindi FROM upasarga_prefix_meta WHERE dhatu_id = ?1 AND prefix = ?2 LIMIT 1"
		err2 := db.QueryRow(q2, dhatuID, upasarga).Scan(&meaning)
		if err2 == nil && meaning.Valid && meaning.String != "" {
			return meaning.String
		}
	}
	if hasTable(db, "upasargachandrika") {
		q2 := "SELECT meaning FROM upasargachandrika WHERE dhatu_id = ?1 AND (upasarga_combination = ?2 OR upasarga_combination = REPLACE(?2, ' + ', '+') OR upasarga_combination = REPLACE(?2, '+', ' + ')) LIMIT 1"
		err2 := db.QueryRow(q2, dhatuID, upasarga).Scan(&meaning)
		if err2 == nil && meaning.Valid && meaning.String != "" {
			return meaning.String
		}
	}
	return ""
}

func fetchLiteraryAttestation(db *sql.DB, dhatuID, word string) map[string]string {
	if word == "" || dhatuID == "" {
		return nil
	}
	var book, text sql.NullString
	if hasTable(db, "prayoga") && hasTable(db, "prayoga_examples") {
		// new schema: prayoga + prayoga_examples are separate but old combined; try old columns first
		q := "SELECT book, literature_text FROM prayoga WHERE dhatu_id = ?1 AND form = ?2 LIMIT 1"
		err := db.QueryRow(q, dhatuID, word).Scan(&book, &text)
		if err == nil && (book.Valid || text.Valid) {
			return map[string]string{
				"book": book.String,
				"text": text.String,
			}
		}
		// new schemaFallback: join prayoga + prayoga_examples
		q2 := `SELECT pe.book, pe.text FROM prayoga p JOIN prayoga_examples pe ON pe.prayoga_id = p.id WHERE p.dhatu_id = ?1 AND p.form = ?2 LIMIT 1`
		err2 := db.QueryRow(q2, dhatuID, word).Scan(&book, &text)
		if err2 == nil && (book.Valid || text.Valid) {
			return map[string]string{
				"book": book.String,
				"text": text.String,
			}
		}
		return nil
	}
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

func fetchVerbsNew(db *sql.DB, word string) []map[string]any {
	var results []map[string]any
	// new schema: conjugation_forms
	query := `SELECT cf.dhatu_id, cf.prefix, cf.category, cf.form_type, cf.form_value FROM conjugation_forms cf WHERE cf.form_value = ?1`
	rows, err := db.Query(query, word)
	if err != nil {
		return results
	}
	defer rows.Close()
	seenVerb := make(map[string]bool)
	for rows.Next() {
		var dhatuID, prefix, category, formType, formValue sql.NullString
		if err := rows.Scan(&dhatuID, &prefix, &category, &formType, &formValue); err != nil {
			continue
		}
		key := dhatuID.String + "|" + prefix.String + "|" + formType.String
		if seenVerb[key] {
			continue
		}
		seenVerb[key] = true
		var root, meaning sql.NullString
		db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('English Meaning','hindI arTa') LIMIT 1", dhatuID.String).Scan(&meaning)
		if !meaning.Valid || meaning.String == "" {
			db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name='arTaH' LIMIT 1", dhatuID.String).Scan(&meaning)
		}
		db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('OpadeSikasvarUpam','mUlaDAtuH','DAtuH','dhatu') LIMIT 1", dhatuID.String).Scan(&root)
		root.String = cleanRoot(root.String)
		// derive lakara/voice from form_type
		ft := formType.String
		voice := "parasmEpadam"
		lakara := ft
		if ft != "" {
			if ft[0] == 'a' || ft[0] == 'A' {
				voice = "Atmanepadam"
				lakara = ft[1:]
			} else if ft[0] == 'p' || ft[0] == 'P' {
				voice = "parasmEpadam"
				lakara = ft[1:]
			}
		}
		derivative := category.String
		if derivative == "ting" {
			derivative = "base"
		}
		prefixedMeaning := fetchUpasargaMeaning(db, dhatuID.String, prefix.String)
		literaryExample := fetchLiteraryAttestation(db, dhatuID.String, word)
		res := map[string]any{
			"type":       "verb",
			"dhatu_id":   dhatuID.String,
			"root":       root.String,
			"meaning":    meaning.String,
			"upasarga":   prefix.String,
			"derivative": derivative,
			"prayoga":    "kartari",
			"lakara":     lakara,
			"voice":      voice,
			"purusha":    "",
			"vacana":     "eka",
		}
		if prefixedMeaning != "" {
			res["prefixed_meaning"] = prefixedMeaning
		}
		if literaryExample != nil {
			res["literature_attestation"] = literaryExample
		}
		res["form_type"] = ft
		results = append(results, res)
	}
	return results
}

func fetchVerbs(db *sql.DB, word string) []map[string]any {
	if hasTable(db, "conjugation_forms") && !hasTable(db, "conjugations") {
		return fetchVerbsNew(db, word)
	}
	var results []map[string]any

	query := `SELECT c.dhatu_id, c.upasarga, c.derivative, c.prayoga, c.lakara, c.voice, c.purusha, c.eka, c.dvi, c.bahu FROM conjugations c WHERE c.eka = ?1 OR c.eka LIKE ?1 || ',%' OR c.eka LIKE '%,' || ?1 OR c.eka LIKE '%,' || ?1 || ',%' OR c.dvi = ?1 OR c.dvi LIKE ?1 || ',%' OR c.dvi LIKE '%,' || ?1 OR c.dvi LIKE '%,' || ?1 || ',%' OR c.bahu = ?1 OR c.bahu LIKE ?1 || ',%' OR c.bahu LIKE '%,' || ?1 OR c.bahu LIKE '%,' || ?1 || ',%'`

	rows, err := db.Query(query, word)
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

		if len(matchedVacanas) == 0 {
			continue
		}

		var root, meaning sql.NullString
		if hasTable(db, "info") {
			db.QueryRow("SELECT value FROM info WHERE dhatu_id=?1 AND key_name IN ('dhatu','OpadeSikasvarUpam','mUlaDAtuH','DAtuH') LIMIT 1", dhatuID.String).Scan(&root)
			root.String = cleanRoot(root.String)
			db.QueryRow("SELECT value FROM info WHERE dhatu_id=?1 AND key_name IN ('english_meaning','hindI_arTa') LIMIT 1", dhatuID.String).Scan(&meaning)
			if !meaning.Valid || meaning.String == "" {
				db.QueryRow("SELECT value FROM info WHERE dhatu_id=?1 AND key_name='arTaH' LIMIT 1", dhatuID.String).Scan(&meaning)
			}
		} else if hasTable(db, "dhatu_info") {
			db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('OpadeSikasvarUpam','mUlaDAtuH','DAtuH','dhatu') LIMIT 1", dhatuID.String).Scan(&root)
			root.String = cleanRoot(root.String)
			db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('English Meaning','hindI arTa') LIMIT 1", dhatuID.String).Scan(&meaning)
			if !meaning.Valid || meaning.String == "" {
				db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name='arTaH' LIMIT 1", dhatuID.String).Scan(&meaning)
			}
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
	// exact match only — no dynamic upasarga sandhi fallback
	return results
}

func fetchParticiplesNew(db *sql.DB, word string) []map[string]any {
	var results []map[string]any
	// exact match only on base/m/f/n
	query := `SELECT dhatu_id, prefix, category, variant, base, m, f, n FROM participle_forms WHERE base = ?1 OR m = ?1 OR f = ?1 OR n = ?1 LIMIT 10`
	rows, err := db.Query(query, word)
	if err != nil {
		return results
	}
	defer rows.Close()
	avyayaPratyayas := []string{"tumun", "ktvA", "lyap", "Ramul"}
	for rows.Next() {
		var dhatuID, prefix, category, variant, base, m, f, n sql.NullString
		if err := rows.Scan(&dhatuID, &prefix, &category, &variant, &base, &m, &f, &n); err != nil {
			continue
		}
		// validity: lyap requires upasarga
		if variant.String == "lyap" && prefix.String == "" {
			continue
		}
		var matchedCols []string
		if base.String == word {
			matchedCols = append(matchedCols, "base_form")
		}
		if m.String == word {
			matchedCols = append(matchedCols, "masculine")
		}
		if f.String == word {
			matchedCols = append(matchedCols, "feminine")
		}
		if n.String == word {
			matchedCols = append(matchedCols, "neuter")
		}
		if len(matchedCols) == 0 {
			continue
		}
		isAvyaya := false
		for _, av := range avyayaPratyayas {
			if variant.String == av {
				isAvyaya = true
				break
			}
		}
		pType := "participle"
		if isAvyaya {
			pType = "avyaya"
		}
		prefixedMeaning := fetchUpasargaMeaning(db, dhatuID.String, prefix.String)
		derivative := category.String
		if derivative == "krut" || derivative == "krt" {
			derivative = "base"
		}
		// fetch root for display (clean anubandha), meaning removed for participles per request
		var rootVal sql.NullString
		if hasTable(db, "dhatu_info") {
			db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('OpadeSikasvarUpam','mUlaDAtuH','DAtuH','dhatu') LIMIT 1", dhatuID.String).Scan(&rootVal)
			rootVal.String = cleanRoot(rootVal.String)
		}
		for _, col := range matchedCols {
			pJSON := map[string]any{
				"type":       pType,
				"dhatu_id":   dhatuID.String,
				"root":       rootVal.String,
				"upasarga":   prefix.String,
				"derivative": derivative,
				"pratyaya":   variant.String,
				"base_form":  base.String,
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

func fetchParticiples(db *sql.DB, word string) []map[string]any {
	if hasTable(db, "participle_forms") && !hasTable(db, "participles") {
		return fetchParticiplesNew(db, word)
	}
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
		var rootVal sql.NullString
		if hasTable(db, "info") {
			db.QueryRow("SELECT value FROM info WHERE dhatu_id=?1 AND key_name IN ('dhatu','OpadeSikasvarUpam','mUlaDAtuH','DAtuH') LIMIT 1", dhatuID.String).Scan(&rootVal)
		} else if hasTable(db, "dhatu_info") {
			db.QueryRow("SELECT value FROM dhatu_info WHERE dhatu_id=?1 AND name IN ('OpadeSikasvarUpam','mUlaDAtuH','DAtuH','dhatu') LIMIT 1", dhatuID.String).Scan(&rootVal)
		}
		rootVal.String = cleanRoot(rootVal.String)

		for _, col := range matchedCols {
			pJSON := map[string]any{
				"type":       pType,
				"dhatu_id":   dhatuID.String,
				"root":       rootVal.String,
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
	// exact match only — no dynamic upasarga sandhi fallback
	return results
}

func analyzeDeclensionNew(db *sql.DB, word string) []map[string]any {
	guessedStems := GetStems(word)
	var results []map[string]any
	seen := make(map[string]bool)
	for _, guess := range guessedStems {
		stem := guess.Stem
		query := "SELECT dhatu_id, prefix, variant, base FROM participle_forms WHERE base = ?1"
		rows, err := db.Query(query, stem)
		var exactMatches []map[string]string
		if err == nil {
			for rows.Next() {
				var dhatuID, prefix, variant, base sql.NullString
				if err := rows.Scan(&dhatuID, &prefix, &variant, &base); err == nil {
					if base.String == stem {
						exactMatches = append(exactMatches, map[string]string{
							"dhatu_id": dhatuID.String,
							"upasarga": prefix.String,
							"pratyaya": variant.String,
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
				dQuery := "SELECT dhatu_id, variant, base FROM participle_forms WHERE base = ?1 AND prefix = ''"
				dRows, err := db.Query(dQuery, strippedStem)
				if err == nil {
					for dRows.Next() {
						var dhatuID, variant, base sql.NullString
						if err := dRows.Scan(&dhatuID, &variant, &base); err == nil {
							if base.String == strippedStem {
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
									"pratyaya":  variant.String,
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

func analyzeDeclension(db *sql.DB, word string) []map[string]any {
	if hasTable(db, "participle_forms") && !hasTable(db, "participles") {
		return analyzeDeclensionNew(db, word)
	}
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
