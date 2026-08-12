package engine

import (
	"database/sql"
	"strings"
	"sync"
)

type irregularEntry struct {
	baseForm string
	gender   string
	cCase    string
	vacana   string
	form     string
}

var (
	irregularCache []irregularEntry
	irregularOnce  sync.Once
)

var irregularCasesList = []string{
	"prathama", "dvitiya", "tritiya", "caturthi",
	"panchami", "sasthi", "saptami", "sambodhana",
}

func loadIrregulars(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM irregulars")
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

		for _, caseName := range irregularCasesList {
			caseVal := rowMap[caseName]
			if caseVal == "" {
				continue
			}

			forms := strings.Split(caseVal, ",")
			for i, form := range forms {
				form = strings.TrimSpace(form)
				if form != "" && i < len(vacanaList) {
					irregularCache = append(irregularCache, irregularEntry{
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

func AnalyzeIrregular(db *sql.DB, word string) []map[string]any {
	irregularOnce.Do(func() {
		loadIrregulars(db)
	})

	var results []map[string]any
	for _, entry := range irregularCache {
		if entry.form == word {
			results = append(results, map[string]any{
				"type":      "irregular_noun",
				"base_form": entry.baseForm,
				"gender":    entry.gender,
				"case":      entry.cCase,
				"vacana":    entry.vacana,
			})
		}
	}

	return results
}