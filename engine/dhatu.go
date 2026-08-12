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
		       COALESCE(i1.value, '') as root,
		       COALESCE(i2.value, '') as meaning,
		       COALESCE(i3.value, '') as gana
		FROM participles p
		LEFT JOIN info i1 ON i1.dhatu_id = p.dhatu_id AND i1.key_name IN ('mUlaDAtuH', 'DAtuH')
		LEFT JOIN info i2 ON i2.dhatu_id = p.dhatu_id AND i2.key_name IN ('aTfaH', 'meaning', 'eng', 'hin')
		LEFT JOIN info i3 ON i3.dhatu_id = p.dhatu_id AND i3.key_name IN ('gaRaH', 'gana')
		WHERE p.base_form LIKE ?1 OR i1.value LIKE ?1 OR i2.value LIKE ?1
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