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

func GenerateVerb(db *sql.DB, root, upasarga, lakara, purusha, voice, prayoga, derivative string) map[string]any {
	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
	} else {
		qExact = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
	}

	var eka, dvi, bahu sql.NullString
	err := db.QueryRow(qExact, root, upasarga, lakara, purusha, voice, prayoga, derivative).Scan(&eka, &dvi, &bahu)
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
			qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
		} else {
			qDyn = "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
		}

		err := db.QueryRow(qDyn, root, lakara, purusha, voice, prayoga, derivative).Scan(&eka, &dvi, &bahu)
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
	isID := isASCIIDigit(root)

	var qExact string
	if isID {
		qExact = "SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
	} else {
		qExact = "SELECT base_form FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
	}

	var baseForm sql.NullString
	err := db.QueryRow(qExact, root, upasarga, pratyaya, derivative).Scan(&baseForm)

	var base string
	found := false

	if err == nil {
		parts := strings.Split(baseForm.String, ",")
		if len(parts) > 0 {
			base = parts[0]
			found = true
		}
	}

	if !found && upasarga != "" {
		var qDyn string
		if isID {
			qDyn = "SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
		} else {
			qDyn = "SELECT base_form FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
		}

		err := db.QueryRow(qDyn, root, pratyaya, derivative).Scan(&baseForm)
		if err == nil {
			parts := strings.Split(baseForm.String, ",")
			if len(parts) > 0 {
				base = ApplyUpasargaSandhi(upasarga, parts[0])
				found = true
			}
		}
	}

	if found {
		avyayas := []string{"tumun", "ktvA", "lyap", "Ramul"}
		for _, av := range avyayas {
			if pratyaya == av {
				return map[string]any{
					"base_form": base,
					"type":      "avyaya",
				}
			}
		}

		declensions, err := DeclineNoun(base, gender)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{
			"base_form":   base,
			"declensions": declensions,
		}
	}

	return map[string]any{
		"error": "Participle combination not found.",
	}
}

func GenerateDeclension(base, gender string) map[string]any {
	declensions, err := DeclineNoun(base, gender)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{
		"base_form":   base,
		"declensions": declensions,
	}
}