package engine

import (
	"database/sql"
	"strings"
	"sync"
)

type numeralEntry struct {
	baseForm string
	gender   string
	cCase    string
	vacana   string
	form     string
}

var (
	numeralCache []numeralEntry
	numeralOnce  sync.Once
)

var numeralCasesList = []string{
	"prathama", "dvitiya", "tritiya", "caturthi", "panchami", "sasthi", "saptami",
}

func loadNumerals(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM numerals")
	if err != nil {
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return
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

		for _, caseName := range numeralCasesList {
			caseVal := rowMap[caseName]
			if caseVal == "" {
				continue
			}

			forms := strings.Split(caseVal, ",")
			for i, form := range forms {
				form = strings.TrimSpace(form)
				if form != "" && i < len(vacanaList) {
					numeralCache = append(numeralCache, numeralEntry{
						baseForm: baseForm,
						gender:   gender,
						cCase:    caseName,
						vacana:   vacanaList[i],
						form:     form,
					})
				}
			}
		}
	}
}

func AnalyzeNumeral(db *sql.DB, word string) []map[string]any {
	numeralOnce.Do(func() {
		loadNumerals(db)
	})

	var results []map[string]any
	for _, entry := range numeralCache {
		if entry.form == word {
			results = append(results, map[string]any{
				"type":      "numeral",
				"base_form": entry.baseForm,
				"gender":    entry.gender,
				"case":      entry.cCase,
				"vacana":    entry.vacana,
			})
		}
	}

	return results
}