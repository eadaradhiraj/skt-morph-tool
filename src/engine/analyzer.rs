use deadpool_sqlite::Pool;
use serde_json::{json, Value};
use super::sandhi::get_upasarga_splits;
use super::stemmer::get_stems;
use super::pronouns::analyze_pronoun;
use super::numerals::analyze_numeral;
use super::irregulars::analyze_irregular;
use super::namadhatu::analyze_namadhatu;

fn strip_accents(word: &str) -> String { word.replace(&['/', '\\', '^'][..], "") }
fn exact_match(row_val: &str, word: &str) -> bool {
    let clean_word = strip_accents(word);
    row_val.split(',').map(|s| strip_accents(s.trim())).any(|s| s == clean_word)
}

fn is_valid_participle(conn: &rusqlite::Connection, dhatu_id: &str, upasarga: &str, pratyaya: &str) -> bool {
    if pratyaya == "lyap" && upasarga.is_empty() { return false; }
    let atmanepada_pratyayas = ["SAnac", "cAnaS", "sya-SAnac", "BAvakarma-SAnac", "sya-BAvakarma-SAnac"];
    if atmanepada_pratyayas.contains(&pratyaya) {
        let mut stmt = conn.prepare("SELECT 1 FROM conjugations WHERE dhatu_id = ?1 AND upasarga = ?2 AND voice = 'Atmanepadam' LIMIT 1").unwrap();
        let mut rows = stmt.query([dhatu_id, upasarga]).unwrap();
        if let Ok(None) = rows.next() { return false; }
    }
    true
}

// ================== VERBS ==================
async fn fetch_verbs(word: String, page: i64, pool: &Pool) -> Vec<Value> {
    let word_clone = word.clone();
    let offset = (page - 1) * 50;
    pool.get().await.unwrap().interact(move |conn| {
        let mut res = Vec::new();
        let search_term = format!("%{}%", strip_accents(&word_clone).chars().map(|c| format!("{}%", c)).collect::<String>());
        
        let mut stmt = conn.prepare("SELECT dhatu_id, upasarga, derivative, prayoga, lakara, voice, purusha, eka, dvi, bahu FROM conjugations WHERE eka LIKE ?1 OR dvi LIKE ?1 OR bahu LIKE ?1 LIMIT 50 OFFSET ?2").unwrap();
        let mut rows = stmt.query(rusqlite::params![search_term, offset]).unwrap();
        while let Ok(Some(row)) = rows.next() {
            let eka: String = row.get("eka").unwrap_or_default();
            let dvi: String = row.get("dvi").unwrap_or_default();
            let bahu: String = row.get("bahu").unwrap_or_default();
            let mut matched_vacanas = Vec::new();
            if exact_match(&eka, &word_clone) { matched_vacanas.push("eka"); }
            if exact_match(&dvi, &word_clone) { matched_vacanas.push("dvi"); }
            if exact_match(&bahu, &word_clone) { matched_vacanas.push("bahu"); }

            for vacana in matched_vacanas {
                res.push(json!({"type": "verb", "dhatu_id": row.get::<&str, String>("dhatu_id").unwrap_or_default(), "upasarga": row.get::<&str, String>("upasarga").unwrap_or_default(), "derivative": row.get::<&str, String>("derivative").unwrap_or_default(), "prayoga": row.get::<&str, String>("prayoga").unwrap_or_default(), "lakara": row.get::<&str, String>("lakara").unwrap_or_default(), "voice": row.get::<&str, String>("voice").unwrap_or_default(), "purusha": row.get::<&str, String>("purusha").unwrap_or_default(), "vacana": vacana}));
            }
        }
        res
    }).await.unwrap()
}

async fn analyze_verb(word: &str, page: i64, pool: &Pool) -> Vec<Value> {
    let mut results = fetch_verbs(word.to_string(), page, pool).await;
    if results.is_empty() && page == 1 {
        for (upa, stripped) in get_upasarga_splits(word) {
            let mut sub_results = fetch_verbs(stripped, 1, pool).await;
            for res in sub_results.iter_mut() {
                if res["upasarga"] == "" {
                    res["upasarga"] = json!(upa.clone());
                    res["note"] = json!("Dynamically matched via Sandhi split");
                    results.push(res.clone());
                }
            }
            if !results.is_empty() { break; }
        }
    }
    results
}

// ================== PARTICIPLES ==================
async fn fetch_participles(word: String, page: i64, pool: &Pool) -> Vec<Value> {
    let word_clone = word.clone();
    let offset = (page - 1) * 50;
    pool.get().await.unwrap().interact(move |conn| {
        let mut res = Vec::new();
        let search_term = format!("%{}%", strip_accents(&word_clone).chars().map(|c| format!("{}%", c)).collect::<String>());
        
        let mut stmt = conn.prepare("SELECT dhatu_id, upasarga, derivative, pratyaya, base_form, masculine, feminine, neuter FROM participles WHERE base_form LIKE ?1 OR masculine LIKE ?1 OR feminine LIKE ?1 OR neuter LIKE ?1 LIMIT 50 OFFSET ?2").unwrap();
        let mut rows = stmt.query(rusqlite::params![search_term, offset]).unwrap();
        
        let avyaya_pratyayas = ["tumun", "ktvA", "lyap", "Ramul"];
        
        while let Ok(Some(row)) = rows.next() {
            let dhatu_id: String = row.get("dhatu_id").unwrap_or_default();
            let upasarga: String = row.get("upasarga").unwrap_or_default();
            let pratyaya: String = row.get("pratyaya").unwrap_or_default();
            
            if !is_valid_participle(&conn, &dhatu_id, &upasarga, &pratyaya) { continue; }
            
            let base_form: String = row.get("base_form").unwrap_or_default();
            let masc: String = row.get("masculine").unwrap_or_default();
            let fem: String = row.get("feminine").unwrap_or_default();
            let neut: String = row.get("neuter").unwrap_or_default();
            
            let mut matched_cols = Vec::new();
            if exact_match(&base_form, &word_clone) { matched_cols.push("base_form"); }
            if exact_match(&masc, &word_clone) { matched_cols.push("masculine"); }
            if exact_match(&fem, &word_clone) { matched_cols.push("feminine"); }
            if exact_match(&neut, &word_clone) { matched_cols.push("neuter"); }

            let p_type = if avyaya_pratyayas.contains(&pratyaya.as_str()) { "avyaya" } else { "participle" };

            for col in matched_cols {
                let mut p_json = json!({"type": p_type, "dhatu_id": dhatu_id.clone(), "upasarga": upasarga.clone(), "derivative": row.get::<&str, String>("derivative").unwrap_or_default(), "pratyaya": pratyaya.clone(), "base_form": base_form.clone()});
                if col == "masculine" || col == "feminine" || col == "neuter" {
                    p_json["gender"] = json!(col); p_json["case"] = json!("prathama"); p_json["vacana"] = json!("eka");
                } else {
                    p_json["note"] = json!("Matched uninflected base form");
                }
                res.push(p_json);
            }
        }
        res
    }).await.unwrap()
}

async fn analyze_participle(word: &str, page: i64, pool: &Pool) -> Vec<Value> {
    let mut results = fetch_participles(word.to_string(), page, pool).await;
    if results.is_empty() && page == 1 {
        for (upa, stripped) in get_upasarga_splits(word) {
            let mut sub_results = fetch_participles(stripped, 1, pool).await;
            for res in sub_results.iter_mut() {
                if res["upasarga"] == "" {
                    res["upasarga"] = json!(upa.clone());
                    res["note"] = json!("Dynamically matched via Sandhi split");
                    results.push(res.clone());
                }
            }
            if !results.is_empty() { break; }
        }
    }
    results
}

// ================== DECLENSIONS (NOUNS) ==================
async fn analyze_declension(word: &str, page: i64, pool: &Pool) -> Vec<Value> {
    let guessed_stems = get_stems(word);
    let guesses_json: Vec<Value> = guessed_stems.into_iter().map(|g| json!({"stem": g.stem, "gender": g.gender, "case": g.case, "vacana": g.vacana})).collect();
    let offset = (page - 1) * 50;
    
    pool.get().await.unwrap().interact(move |conn| {
        let mut results = Vec::new();
        for guess in guesses_json {
            let stem = guess["stem"].as_str().unwrap();
            let search_term = format!("%{}%", stem);
            
            let mut stmt = conn.prepare("SELECT dhatu_id, upasarga, pratyaya, base_form FROM participles WHERE base_form LIKE ?1 LIMIT 50 OFFSET ?2").unwrap();
            let mut rows = stmt.query(rusqlite::params![search_term, offset]).unwrap();
            
            let mut exact_matches = Vec::new();
            while let Ok(Some(row)) = rows.next() {
                let base_form: String = row.get("base_form").unwrap_or_default();
                if exact_match(&base_form, stem) {
                    exact_matches.push(json!({"dhatu_id": row.get::<&str, String>("dhatu_id").unwrap_or_default(), "upasarga": row.get::<&str, String>("upasarga").unwrap_or_default(), "pratyaya": row.get::<&str, String>("pratyaya").unwrap_or_default()}));
                }
            }
            
            if !exact_matches.is_empty() {
                for m in exact_matches {
                    let mut r = guess.clone();
                    r["type"] = json!("declension"); r["base_form"] = json!(stem);
                    r["dhatu_id"] = m["dhatu_id"].clone(); r["upasarga"] = m["upasarga"].clone(); r["pratyaya"] = m["pratyaya"].clone();
                    results.push(r);
                }
            } else if page == 1 {
                let mut found_dynamic = false;
                for (upa, stripped_stem) in get_upasarga_splits(stem) {
                    let s_term = format!("%{}%", stripped_stem);
                    let mut d_stmt = conn.prepare("SELECT dhatu_id, pratyaya, base_form FROM participles WHERE base_form LIKE ?1 AND upasarga = '' LIMIT 50").unwrap();
                    let mut d_rows = d_stmt.query([&s_term]).unwrap();
                    
                    while let Ok(Some(row)) = d_rows.next() {
                        let base_form: String = row.get("base_form").unwrap_or_default();
                        if exact_match(&base_form, &stripped_stem) {
                            let mut r = guess.clone();
                            r["type"] = json!("declension"); r["base_form"] = json!(stem);
                            r["dhatu_id"] = json!(row.get::<&str, String>("dhatu_id").unwrap_or_default());
                            r["upasarga"] = json!(upa.clone());
                            r["pratyaya"] = json!(row.get::<&str, String>("pratyaya").unwrap_or_default());
                            r["note"] = json!("Dynamic Upasarga Match");
                            results.push(r); found_dynamic = true;
                        }
                    }
                    if found_dynamic { break; }
                }
                if !found_dynamic {
                    let mut r = guess.clone();
                    r["type"] = json!("declension"); r["base_form"] = json!(stem);
                    r["dhatu_id"] = Value::Null; r["upasarga"] = Value::Null; r["pratyaya"] = Value::Null;
                    results.push(r);
                }
            }
        }
        results
    }).await.unwrap()
}

pub async fn analyze(word: &str, page: i64, pool: &Pool) -> Value {
    let verbs = analyze_verb(word, page, pool).await;
    let participles = analyze_participle(word, page, pool).await;
    let declensions = analyze_declension(word, page, pool).await;
    
    // We don't paginate irregulars, numerals, or pronouns as they are tiny static lists
    let (pronouns, numerals, irregulars, namadhatus) = if page == 1 {
        (analyze_pronoun(word, pool).await, analyze_numeral(word, pool).await, analyze_irregular(word, pool).await, analyze_namadhatu(word))
    } else {
        (vec![], vec![], vec![], vec![])
    };

    let has_more = verbs.len() >= 50 || participles.len() >= 50 || declensions.len() >= 50;

    json!({
        "searched_word": word,
        "page": page,
        "has_more": has_more,
        "verbs": verbs,
        "participles": participles,
        "declensions": declensions,
        "pronouns": pronouns,
        "numerals": numerals,
        "irregulars": irregulars,
        "namadhatus": namadhatus
    })
}
