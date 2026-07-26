use deadpool_sqlite::Pool;
use serde_json::{json, Value};

const CASES: [&str; 7] = ["prathama", "dvitiya", "tritiya", "caturthi", "panchami", "sasthi", "saptami"];

pub async fn analyze_numeral(word: &str, pool: &Pool) -> Vec<Value> {
    let word_clone = word.to_string();
    pool.get().await.unwrap().interact(move |conn| {
        let mut results = Vec::new();
        let mut stmt = conn.prepare("SELECT * FROM numerals").unwrap();
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
                            "type": "numeral",
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
