package engine

import (
	"database/sql"
	"strings"
)

var numeralCases = []string{
	"prathama", "dvitiya", "tritiya", "caturthi", "panchami", "sasthi", "saptami",
}

// AnalyzeNumeral looks up inflected numeral forms in the numerals table
func AnalyzeNumeral(db *sql.DB, word string) []map[string]any {
	var results []map[string]any

	rows, err := db.Query("SELECT * FROM numerals")
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

		for _, caseName := range numeralCases {
			caseVal := rowMap[caseName]
			if caseVal == "" {
				continue
			}

			forms := strings.Split(caseVal, ",")
			for i, form := range forms {
				form = strings.TrimSpace(form)
				if form == word && form != "" && i < len(vacanaList) {
					results = append(results, map[string]any{
						"type":      "numeral",
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