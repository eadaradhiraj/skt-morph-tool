package engine

import (
	"database/sql"
	"strings"
)

var irregularCases = []string{
	"prathama", "dvitiya", "tritiya", "caturthi",
	"panchami", "sasthi", "saptami", "sambodhana",
}

func AnalyzeIrregular(db *sql.DB, word string) []map[string]any {
	var results []map[string]any

	rows, err := db.Query("SELECT * FROM irregulars")
	if err != nil {
		return results
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return results
	}

	vacanaList := []string{"eka", "dvi", "bahu"}

	for rows.Next() {
		// Scan row into a dynamic map
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		rowMap := make(map[string]string)
		for i, colName := range cols {
			if val, ok := columns[i].([]byte); ok {
				rowMap[colName] = string(val)
			} else if val, ok := columns[i].(string); ok {
				rowMap[colName] = val
			}
		}

		baseForm := rowMap["base_form"]
		gender := rowMap["gender"]

		for _, caseName := range irregularCases {
			caseVal := rowMap[caseName]
			if caseVal == "" {
				continue
			}

			forms := strings.Split(caseVal, ",")
			for i, form := range forms {
				form = strings.TrimSpace(form)
				if form == word && form != "" && i < len(vacanaList) {
					results = append(results, map[string]any{
						"type":      "irregular_noun",
						"base_form": baseForm,
						"gender":    gender,
						"case":      caseName,
						"vacana":    vacanaList[i],
					})
				}
			}
		}
	}

	return results
}

func declineIrregular(base string, gender string) (map[string][]string, bool) {
	if base == "go" && (gender == "masculine" || gender == "feminine") {
		return map[string][]string{
			"prathama": {"gOH", "gAvO", "gAvaH"},
			"dvitiya":  {"gAm", "gAvO", "gAH"},
			"tritiya":  {"gavA", "goByAm", "goBiH"},
		}, true
	}
	if base == "strI" && gender == "feminine" {
		return map[string][]string{
			"prathama": {"strI", "striyO", "striyaH"},
			"dvitiya":  {"striyam", "striyO", "striyaH"},
		}, true
	}
	return nil, false
}