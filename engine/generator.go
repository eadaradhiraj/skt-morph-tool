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
	if strings.Contains(vLower, "atman") {
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
	case "law", "lat", "laT":
		return "laT"
	case "low", "lot", "loT":
		return "loT"
	case "lan", "laN":
		return "laN"
	case "vidilin", "vidhilin", "vidhiliN", "vidiliN":
		return "vidhiliN"
	case "lfw", "lrt", "lfT":
		return "lfT"
	case "liw", "lit", "liT":
		return "liT"
	case "luw", "lut", "luT":
		return "luT"
	case "asirling", "asirliN", "asir-lin", "ASIrliN":
		return "ASIrliN"
	case "lun", "luN":
		return "luN"
	case "lrn", "lfwn", "lrN":
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
	case "sotf", "satf", "Satf":
		return "Satf"
	case "sanac", "SAnac":
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

func hasTableGen(db *sql.DB, name string) bool {
	var n string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?1", name).Scan(&n)
	return err == nil && n == name
}

func lakaraSuffix(lakara string) string {
	switch strings.ToLower(lakara) {
	case "lat", "law":
		return "lat"
	case "lan", "laN":
		return "lang"
	case "lit", "liw":
		return "lit"
	case "lot", "low":
		return "lot"
	case "lun", "luN":
		return "lung"
	case "lut", "luw":
		return "lut"
	case "lrt", "lfw", "lft":
		return "lrut"
	case "lrn", "lfn", "lrv":
		return "lrung"
	case "vidhilin", "vidhiliN":
		return "vidhiling"
	case "asirlin", "asirling", "asirliN", "ASIrliN":
		return "ashirling"
	default:
		return strings.ToLower(lakara)
	}
}

func formTypeFor(lakara, voice string) string {
	suffix := lakaraSuffix(lakara)
	if voice == "Atmanepadam" {
		return "a" + suffix
	}
	return "p" + suffix
}

func ekaDviBahuFromForms(forms []string, purusha string) (string, string, string) {
	n := len(forms)
	if n == 0 {
		return "", "", ""
	}
	p := strings.ToLower(purusha)
	if p == "prathama" || p == "praTama" || p == "pratama" {
		if n == 10 { // pashirling/vidhiling etc: eka has 2 forms
			return forms[0] + ", " + forms[1], forms[2], forms[3]
		}
		if n == 13 { // lot
			return forms[0] + ", " + forms[1] + ", " + forms[2], forms[3], forms[4]
		}
		if n >= 3 {
			return forms[0], forms[1], forms[2]
		}
	} else if p == "madhyama" || p == "madyama" {
		if n == 10 {
			return forms[4], forms[5], forms[6]
		}
		if n == 13 {
			return forms[5] + ", " + forms[6] + ", " + forms[7], forms[8], forms[9]
		}
		if n >= 6 {
			return forms[3], forms[4], forms[5]
		}
	} else if p == "uttama" {
		if n == 10 {
			return forms[7], forms[8], forms[9]
		}
		if n == 13 {
			return forms[10], forms[11], forms[12]
		}
		if n >= 9 {
			return forms[6], forms[7], forms[8]
		}
	}
	// fallback generic 9
	if n >= 9 {
		if p == "prathama" {
			return forms[0], forms[1], forms[2]
		}
		if p == "madhyama" {
			return forms[3], forms[4], forms[5]
		}
		return forms[6], forms[7], forms[8]
	}
	return forms[0], "", ""
}

func GenerateVerbNew(db *sql.DB, root, upasarga, lakara, purusha, voice, prayoga, derivative string) map[string]any {
	// new schema: conjugation_forms
	isID := isASCIIDigit(root)
	var dhatuID string
	if isID {
		dhatuID = root
	} else {
		// lookup dhatu_id via dhatu_info value
		err := db.QueryRow("SELECT dhatu_id FROM dhatu_info WHERE value=?1 LIMIT 1", root).Scan(&dhatuID)
		if err != nil {
			// try like
			db.QueryRow("SELECT dhatu_id FROM dhatu_info WHERE value LIKE ?1 || '%' LIMIT 1", root).Scan(&dhatuID)
		}
		if dhatuID == "" {
			return map[string]any{"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible."}
		}
	}
	// derivative: primary/base -> ting
	cat := derivative
	if cat == "mUla" || cat == "base" {
		cat = "ting"
	} else if cat == "" {
		cat = "ting"
	}
	ft := formTypeFor(lakara, voice)
	// query all forms for this combination ordered by position
	rows, err := db.Query("SELECT form_value FROM conjugation_forms WHERE dhatu_id=?1 AND prefix=?2 AND category=?3 AND form_type=?4 ORDER BY position", dhatuID, upasarga, cat, ft)
	if err != nil {
		return map[string]any{"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible."}
	}
	defer rows.Close()
	var forms []string
	for rows.Next() {
		var fv sql.NullString
		rows.Scan(&fv)
		forms = append(forms, fv.String)
	}
	if len(forms) == 0 {
		return map[string]any{"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible."}
	}
	eka, dvi, bahu := ekaDviBahuFromForms(forms, purusha)
	if eka == "" && dvi == "" && bahu == "" {
		return map[string]any{"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible."}
	}
	return map[string]any{"eka": eka, "dvi": dvi, "bahu": bahu}
}

func GenerateVerb(db *sql.DB, root, upasarga, lakara, purusha, voice, prayoga, derivative string) map[string]any {
	if hasTableGen(db, "conjugation_forms") && !hasTableGen(db, "conjugations") {
		// normalize for new path as well
		root = strings.TrimSpace(root)
		upasarga = strings.TrimSpace(upasarga)
		lakara = normalizeLakara(lakara)
		purusha = normalizePurusha(purusha)
		voice = normalizeVoice(voice)
		prayoga = normalizePrayoga(prayoga)
		derivative = normalizeDerivative(derivative)
		return GenerateVerbNew(db, root, upasarga, lakara, purusha, voice, prayoga, derivative)
	}
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

	return map[string]any{
		"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible.",
	}
}

func GenerateParticipleNew(db *sql.DB, root, upasarga, pratyaya, gender, derivative string) map[string]any {
	isID := isASCIIDigit(root)
	var dhatuID string
	if isID {
		dhatuID = root
	} else {
		db.QueryRow("SELECT dhatu_id FROM dhatu_info WHERE value=?1 LIMIT 1", root).Scan(&dhatuID)
		if dhatuID == "" {
			db.QueryRow("SELECT dhatu_id FROM dhatu_info WHERE value LIKE ?1 || '%' LIMIT 1", root).Scan(&dhatuID)
		}
		if dhatuID == "" {
			return map[string]any{"error": "Participle combination not found."}
		}
	}
	cat := derivative
	if cat == "mUla" || cat == "base" {
		cat = "krut"
	}
	var baseForm, m, f, n sql.NullString
	err := db.QueryRow("SELECT base, m, f, n FROM participle_forms WHERE dhatu_id=?1 AND prefix=?2 AND variant=?3 AND category=?4 LIMIT 1", dhatuID, upasarga, pratyaya, cat).Scan(&baseForm, &m, &f, &n)
	if err != nil {
		return map[string]any{"error": "Participle combination not found."}
	}
	avyayas := []string{"tumun", "ktvA", "lyap", "Ramul"}
	for _, av := range avyayas {
		if pratyaya == av {
			return map[string]any{"base_form": cleanParticipleStem(baseForm.String), "type": "avyaya"}
		}
	}
	var selectedBase string
	if gender == "feminine" && f.Valid && f.String != "" {
		selectedBase = cleanParticipleStem(f.String)
	} else if gender == "neuter" && n.Valid && n.String != "" {
		selectedBase = cleanParticipleStem(n.String)
	} else if gender == "masculine" && m.Valid && m.String != "" {
		selectedBase = cleanParticipleStem(m.String)
	} else if baseForm.Valid {
		selectedBase = cleanParticipleStem(baseForm.String)
	}
	declensions, err := DeclineNoun(selectedBase, gender)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{"base_form": selectedBase, "declensions": declensions}
}

func GenerateParticiple(db *sql.DB, root, upasarga, pratyaya, gender, derivative string) map[string]any {
	if hasTableGen(db, "participle_forms") && !hasTableGen(db, "participles") {
		root = strings.TrimSpace(root)
		upasarga = strings.TrimSpace(upasarga)
		pratyaya = normalizePratyaya(pratyaya)
		gender = normalizeGender(gender)
		derivative = normalizeDerivative(derivative)
		return GenerateParticipleNew(db, root, upasarga, pratyaya, gender, derivative)
	}
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

	found := err == nil

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
