package engine

import (
	"database/sql"
	"fmt"
)

// SearchDhatu queries Dhatu information and participles matching a search query
func SearchDhatu(db *sql.DB, query string) []map[string]any {
	var results []map[string]any

	searchTerm := fmt.Sprintf("%%%s%%", query)

	sqlQuery := `
		SELECT p.dhatu_id, p.pratyaya, p.base_form, p.upasarga,
		       (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('mUlaDAtuH', 'DAtuH') LIMIT 1) as root,
		       (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('aTfaH', 'meaning', 'eng', 'hin') LIMIT 1) as meaning,
		       (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('gaRaH', 'gana') LIMIT 1) as gana
		FROM participles p
		WHERE p.base_form LIKE ?1 OR root LIKE ?1 OR meaning LIKE ?1
		LIMIT 100
	`

	rows, err := db.Query(sqlQuery, searchTerm)
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		var dhatuID, pratyaya, baseForm, upasarga, root, meaning, gana sql.NullString
		if err := rows.Scan(&dhatuID, &pratyaya, &baseForm, &upasarga, &root, &meaning, &gana); err != nil {
			continue
		}

		results = append(results, map[string]any{
			"dhatu_id":  dhatuID.String,
			"root":      root.String,
			"meaning":   meaning.String,
			"gana":      gana.String,
			"pratyaya":  pratyaya.String,
			"upasarga":  upasarga.String,
			"base_form": baseForm.String,
		})
	}

	return results
}