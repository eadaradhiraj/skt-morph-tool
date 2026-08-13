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
	p = strings.ToLower(strings.TrimSpace(p))
	switch p {
	case "kartari", "active":
		return "Kartari"
	case "karmani", "passive":
		return "Karmani"
	case "bhave", "impersonal":
		return "Bhave"
	default:
		if p == "" {
			return "Kartari"
		}
		return p
	}
}

func normalizeVoice(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	switch v {
	case "parasmaipada", "parasmaipadam":
		return "Parasmaipadam"
	case "atmanepada", "atmanepadam":
		return "Atmanepadam"
	default:
		return v
	}
}

func normalizeDerivative(d string) string {
	d = strings.TrimSpace(d)
	if d == "" {
		return "mUla"
	}
	return d
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
	lakara = strings.TrimSpace(lakara)
	purusha = strings.TrimSpace(purusha)
	voice = normalizeVoice(voice)
	prayoga = normalizePrayoga(prayoga)
	derivative = normalizeDerivative(derivative)

	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		if voice != "" {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
		} else {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND prayoga=?5 AND derivative=?6 ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
		}
	} else {
		if voice != "" {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
		} else {
			qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND prayoga=?5 AND derivative=?6 ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
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
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
			} else {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND lakara=?2 AND purusha=?3 AND prayoga=?4 AND derivative=?5 ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
			}
		} else {
			if voice != "" {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
			} else {
				qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND lakara=?2 AND purusha=?3 AND prayoga=?4 AND derivative=?5 ORDER BY CASE WHEN voice='Parasmaipadam' THEN 1 ELSE 2 END LIMIT 1"
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
	pratyaya = strings.TrimSpace(pratyaya)
	gender = normalizeGender(gender)
	derivative = normalizeDerivative(derivative)

	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		qExact = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id=?1 AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
	} else {
		qExact = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
	}

	var baseForm, masc, fem, neut sql.NullString
	err := db.QueryRow(qExact, root, upasarga, pratyaya, derivative).Scan(&baseForm, &masc, &fem, &neut)

	found := false

	if err == nil {
		found = true
	} else if upasarga != "" {
		var qDyn string
		if isID {
			qDyn = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id=?1 AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
		} else {
			qDyn = "SELECT base_form, masculine, feminine, neuter FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1 OR value LIKE ?1 || '~%' OR value LIKE ?1 || '%') AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
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
