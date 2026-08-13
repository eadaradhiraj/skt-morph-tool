package engine

import (
	"database/sql"
	"strings"
	"unicode"
)

func isASCIIDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsDigit(rune(s[0]))
}

func normalizeGender(g string) string {
	g = strings.ToLower(strings.TrimSpace(g))
	switch g {
	case "m", "masc", "masculine":
		return "masculine"
	case "f", "fem", "feminine":
		return "feminine"
	case "n", "neut", "neuter":
		return "neuter"
	default:
		return g
	}
}

func normalizePrayoga(p string) string {
	p = strings.TrimSpace(p)
	pLower := strings.ToLower(p)
	if strings.Contains(pLower, "kartar") {
		return "Kartari"
	}
	if strings.Contains(pLower, "karm") {
		return "Karmani"
	}
	if strings.Contains(pLower, "bhav") || strings.Contains(pLower, "bav") {
		return "Bhave"
	}
	if p == "" {
		return "Kartari"
	}
	return p
}

func normalizeVoice(v string) string {
	v = strings.TrimSpace(v)
	vLower := strings.ToLower(v)
	if strings.Contains(vLower, "parasm") || strings.Contains(vLower, "active") {
		return "Parasmaipadam"
	}
	if strings.Contains(vLower, "atman") || strings.Contains(vLower, "atman") {
		return "Atmanepadam"
	}
	return v
}

func normalizeDerivative(d string) string {
	d = strings.TrimSpace(d)
	if d == "" || d == "base" || strings.EqualFold(d, "mula") || strings.EqualFold(d, "primary") {
		return "mUla"
	}
	return d
}

func normalizeLakara(l string) string {
	l = strings.TrimSpace(l)
	switch strings.ToLower(l) {
	case "law", "lat", "lat", "lat", "laT":
		return "laT"
	case "low", "lot", "loT":
		return "loT"
	case "lan", "lan", "laN":
		return "laN"
	case "vidilin", "vidhilin", "vidhiliN", "vidiliN":
		return "vidhiliN"
	case "lfw", "lrt", "lfT":
		return "lfT"
	case "liw", "lit", "liT":
		return "liT"
	case "luw", "lut", "luT":
		return "luT"
	case "asirling", "asirlin", "asir-lin", "ASIrliN", "asirliN":
		return "ASIrliN"
	case "lun", "luN":
		return "luN"
	case "lrn", "lfwN", "lrN":
		return "lrN"
	default:
		return l
	}
}

func normalizePurusha(p string) string {
	p = strings.TrimSpace(p)
	pLower := strings.ToLower(p)
	switch pLower {
	case "pratama", "prathama", "3rd", "third":
		return "prathama"
	case "madyama", "madhyama", "2nd", "second":
		return "madhyama"
	case "uttama", "1st", "first":
		return "uttama"
	default:
		return p
	}
}

func normalizePratyaya(pr string) string {
	pr = strings.TrimSpace(pr)
	switch strings.ToLower(pr) {
	case "sotf", "satf", "satf", "Satf":
		return "Satf"
	case "sanac", "sanac", "SAnac":
		return "SAnac"
	case "tfc", "trc":
		return "trc"
	default:
		return pr
	}
}

func cleanParticipleStem(raw string) string {
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, ",")
	first := strings.TrimSpace(parts[0])
	return strings.TrimSuffix(first, "H")
}

func GenerateVerb(db *sql.DB, root, upasarga, lakara, purusha, voice, prayoga, derivative string) map[string]any {
	root = strings.TrimSpace(root)
	upasarga = strings.TrimSpace(upasarga)
	lakara = normalizeLakara(lakara)
	purusha = normalizePurusha(purusha)
	voice = normalizeVoice(voice)
	prayoga = normalizePrayoga(prayoga)
	derivative = normalizeDerivative(derivative)

	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		if voice != "" {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND (lakara=?3 OR lakara=REPLACE(?3,'T','w')) AND (purusha=?4 OR purusha='praTama' OR purusha='maDyama') AND voice=?5 AND prayoga=?6 AND (derivative=?7 OR derivative='base') LIMIT 1"
		} else {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND (lakara=?3 OR lakara=REPLACE(?3,'T','w')) AND (purusha=?4 OR purusha='praTama' OR purusha='maDyama') AND prayoga=?5 AND (derivative=?6 OR derivative='base') ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
		}
	} else {
		if voice != "" {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND (lakara=?3 OR lakara=REPLACE(?3,'T','w')) AND (purusha=?4 OR purusha='praTama' OR purusha='maDyama') AND voice=?5 AND prayoga=?6 AND (derivative=?7 OR derivative='base') LIMIT 1"
		} else {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND (lakara=?3 OR lakara=REPLACE(?3,'T','w')) AND (purusha=?4 OR purusha='praTama' OR purusha='maDyama') AND prayoga=?5 AND (derivative=?6 OR derivative='base') ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
		}
	}

	var eka, dvi, bahu sql.NullString
	var err error
	if voice != "" {
		err = db.QueryRow(qExact, root, upasarga, lakara, purusha, voice, prayoga, derivative).Scan(&eka, &dvi, &bahu)
	} else {
		err = db.QueryRow(qExact, root, upasarga, lakara, purusha, prayoga, derivative).Scan(&eka, &dvi, &bahu)
	}

	if err == nil {
		return map[string]any{
			"eka":  eka.String,
			"dvi":  dvi.String,
			"bahu": bahu.String,
		}
	}

	if upasarga != "" {
		var qDyn string
		if isID {
			if voice != "" {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND (lakara=?2 OR lakara=REPLACE(?2,'T','w')) AND (purusha=?3 OR purusha='praTama' OR purusha='maDyama') AND voice=?4 AND prayoga=?5 AND (derivative=?6 OR derivative='base') LIMIT 1"
			} else {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND (lakara=?2 OR lakara=REPLACE(?2,'T','w')) AND (purusha=?3 OR purusha='praTama' OR purusha='maDyama') AND prayoga=?4 AND (derivative=?5 OR derivative='base') ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
			}
		} else {
			if voice != "" {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND (lakara=?2 OR lakara=REPLACE(?2,'T','w')) AND (purusha=?3 OR purusha='praTama' OR purusha='maDyama') AND voice=?4 AND prayoga=?5 AND (derivative=?6 OR derivative='base') LIMIT 1"
			} else {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND (lakara=?2 OR lakara=REPLACE(?2,'T','w')) AND (purusha=?3 OR purusha='praTama' OR purusha='maDyama') AND prayoga=?4 AND (derivative=?5 OR derivative='base') ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
			}
		}

		if voice != "" {
			err = db.QueryRow(qDyn, root, lakara, purusha, voice, prayoga, derivative).Scan(&eka, &dvi, &bahu)
		} else {
			err = db.QueryRow(qDyn, root, lakara, purusha, prayoga, derivative).Scan(&eka, &dvi, &bahu)
		}

		if err == nil {
			fEka := ApplyUpasargaSandhi(upasarga, eka.String)
			fDvi := ApplyUpasargaSandhi(upasarga, dvi.String)
			fBahu := ApplyUpasargaSandhi(upasarga, bahu.String)
			return map[string]any{
				"eka":  fEka,
				"dvi":  fDvi,
				"bahu": fBahu,
				"note": "Dynamically Sandhi-fused",
			}
		}
	}

	return map[string]any{
		"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible.",
	}
}

func GenerateParticiple(db *sql.DB, root, upasarga, pratyaya, gender, derivative string) map[string]any {
	root = strings.TrimSpace(root)
	upasarga = strings.TrimSpace(upasarga)
	pratyaya = normalizePratyaya(pratyaya)
	gender = normalizeGender(gender)
	derivative = normalizeDerivative(derivative)

	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		qExact = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id=?1 AND upasarga=?2 AND (pratyaya=?3 OR pratyaya='Sotf' OR pratyaya='tfc') AND (derivative=?4 OR derivative='base') LIMIT 1"
	} else {
		qExact = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND (pratyaya=?3 OR pratyaya='Sotf' OR pratyaya='tfc') AND (derivative=?4 OR derivative='base') LIMIT 1"
	}

	var baseForm, masc, fem, neut sql.NullString
	err := db.QueryRow(qExact, root, upasarga, pratyaya, derivative).Scan(&baseForm, &masc, &fem, &neut)

	found := false

	if err == nil {
		found = true
	} else if upasarga != "" {
		var qDyn string
		if isID {
			qDyn = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id=?1 AND upasarga='' AND (pratyaya=?2 OR pratyaya='Sotf' OR pratyaya='tfc') AND (derivative=?3 OR derivative='base') LIMIT 1"
		} else {
			qDyn = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND (pratyaya=?2 OR pratyaya='Sotf' OR pratyaya='tfc') AND (derivative=?3 OR derivative='base') LIMIT 1"
		}

		err := db.QueryRow(qDyn, root, pratyaya, derivative).Scan(&baseForm, &masc, &fem, &neut)
		if err == nil {
			found = true
			if baseForm.Valid {
				baseForm.String = ApplyUpasargaSandhi(upasarga, baseForm.String)
			}
			if fem.Valid {
				fem.String = ApplyUpasargaSandhi(upasarga, fem.String)
			}
			if masc.Valid {
				masc.String = ApplyUpasargaSandhi(upasarga, masc.String)
			}
			if neut.Valid {
				neut.String = ApplyUpasargaSandhi(upasarga, neut.String)
			}
		}
	}

	if found {
		avyayas := []string{"tumun", "ktvA", "lyap", "Ramul"}
		for _, av := range avyayas {
			if pratyaya == av {
				return map[string]any{
					"base_form": cleanParticipleStem(baseForm.String),
					"type":      "avyaya",
				}
			}
		}

		var selectedBase string
		if gender == "feminine" && fem.Valid && fem.String != "" {
			selectedBase = cleanParticipleStem(fem.String)
		} else if gender == "neuter" && neut.Valid && neut.String != "" {
			selectedBase = cleanParticipleStem(neut.String)
		} else if gender == "masculine" && masc.Valid && masc.String != "" {
			selectedBase = cleanParticipleStem(masc.String)
		} else if baseForm.Valid {
			selectedBase = cleanParticipleStem(baseForm.String)
		}

		declensions, err := DeclineNoun(selectedBase, gender)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{
			"base_form":   selectedBase,
			"declensions": declensions,
		}
	}

	return map[string]any{
		"error": "Participle combination not found.",
	}
}

func GenerateDeclension(base, gender string) map[string]any {
	base = strings.TrimSpace(base)
	gender = normalizeGender(gender)

	declensions, err := DeclineNoun(base, gender)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{
		"base_form":   base,
		"declensions": declensions,
	}
}
