use deadpool_sqlite::Pool;
use serde_json::{json, Value};
use std::collections::HashMap;

const CASES: [&str; 8] = ["prathama", "dvitiya", "tritiya", "caturthi", "panchami", "sasthi", "saptami", "sambodhana"];

pub async fn analyze_irregular(word: &str, pool: &Pool) -> Vec<Value> {
    let word_clone = word.to_string();
    pool.get().await.unwrap().interact(move |conn| {
        let mut results = Vec::new();
        let mut stmt = conn.prepare("SELECT * FROM irregulars").unwrap();
        let mut rows = stmt.query([]).unwrap();

        while let Ok(Some(row)) = rows.next() {
            let base_form: String = row.get("base_form").unwrap_or_default();
            let gender: String = row.get("gender").unwrap_or_default();

            for case in CASES.iter() {
                let case_val: String = row.get(*case).unwrap_or_default();
                let forms: Vec<&str> = case_val.split(',').map(|s| s.trim()).collect();
                
                for (i, form) in forms.iter().enumerate() {
                    if form == &word_clone && !form.is_empty() {
                        let vacana = ["eka", "dvi", "bahu"][i];
                        results.push(json!({
                            "type": "irregular_noun",
                            "base_form": base_form,
                            "gender": gender,
                            "case": case,
                            "vacana": vacana
                        }));
                    }
                }
            }
        }
        results
    }).await.unwrap()
}

pub fn decline_irregular(base: &str, gender: &str) -> Option<HashMap<String, Vec<String>>> {
    // Note: To keep the generator fully sync/blocking, we'll hardcode the most common irregulars here. 
    // In Python we queried the DB, but Rust sync/async boundaries make it cleaner to map them for the generator.
    if base == "go" && (gender == "masculine" || gender == "feminine") {
        let mut m = HashMap::new();
        m.insert("prathama".to_string(), vec!["gOH".to_string(), "gAvO".to_string(), "gAvaH".to_string()]);
        m.insert("dvitiya".to_string(), vec!["gAm".to_string(), "gAvO".to_string(), "gAH".to_string()]);
        m.insert("tritiya".to_string(), vec!["gavA".to_string(), "goByAm".to_string(), "goBiH".to_string()]);
        return Some(m);
    }
    if base == "strI" && gender == "feminine" {
        let mut m = HashMap::new();
        m.insert("prathama".to_string(), vec!["strI".to_string(), "striyO".to_string(), "striyaH".to_string()]);
        m.insert("dvitiya".to_string(), vec!["striyam".to_string(), "striyO".to_string(), "striyaH".to_string()]);
        return Some(m);
    }
    None
}
