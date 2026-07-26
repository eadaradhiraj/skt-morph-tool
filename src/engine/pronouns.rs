use deadpool_sqlite::Pool;
use serde_json::{json, Value};

const CASES: [&str; 7] = ["prathama", "dvitiya", "tritiya", "caturthi", "panchami", "sasthi", "saptami"];

pub async fn analyze_pronoun(word: &str, pool: &Pool) -> Vec<Value> {
    let word_clone = word.to_string();
    pool.get().await.unwrap().interact(move |conn| {
        let mut results = Vec::new();
        // Look up all pronouns (we assume Pronouns were added to the DB alongside Irregulars and Numerals)
        // Note: For now we fall back to a hardcoded mapping since Pronouns didn't explicitly get their own DB table in the Python step yet.
        // We will hardcode the most common ones to match the Python script.
        let pronoun_map = vec![
            ("tad", "masculine", "prathama", vec!["saH", "tO", "te"]),
            ("tad", "masculine", "dvitiya", vec!["tam", "tO", "tAn"]),
            ("tad", "masculine", "caturthi", vec!["tasmE", "tAByAm", "teByaH"]),
            ("kim", "masculine", "prathama", vec!["kaH", "kO", "ke"]),
            ("asmad", "any", "prathama", vec!["aham", "AvAm", "vayam"]),
            ("yuzmad", "any", "prathama", vec!["tvam", "yuvAm", "yUyam"]),
        ];

        for (base, gender, case, forms) in pronoun_map {
            for (i, form) in forms.iter().enumerate() {
                if form == &word_clone {
                    let vacana = ["eka", "dvi", "bahu"][i];
                    results.push(json!({
                        "type": "pronoun",
                        "base_form": base,
                        "gender": gender,
                        "case": case,
                        "vacana": vacana
                    }));
                }
            }
        }
        results
    }).await.unwrap()
}
