use deadpool_sqlite::Pool;
use serde_json::{json, Value};
use crate::engine::sandhi::apply_upasarga_sandhi;
use crate::engine::declension::decline_noun;

pub async fn generate_verb(pool: &Pool, dhatu_id: &str, root: &str, upasarga: &str, lakara: &str, purusha: &str, voice: &str) -> Value {
    let d = dhatu_id.to_string(); let r = root.to_string(); let u = upasarga.to_string(); let l = lakara.to_string(); let p = purusha.to_string(); let v = voice.to_string();

    let row_json: Option<Value> = pool.get().await.unwrap().interact(move |conn| {
        let target_id = if !d.is_empty() { d } else if !r.is_empty() {
            let mut stmt = conn.prepare("SELECT dhatu_id FROM info WHERE value = ?1 LIMIT 1").unwrap();
            let mut rows = stmt.query([&r]).unwrap();
            if let Ok(Some(row)) = rows.next() { row.get::<usize, String>(0).unwrap_or_default() } else { return Some(json!({"error": format!("Root '{}' not found.", r)})); }
        } else { return Some(json!({"error": "You must provide either 'dhatu_id' or 'root'."})); };

        // Try exact match
        let mut stmt = conn.prepare("SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND derivative='base' LIMIT 1").unwrap();
        let mut rows = stmt.query([&target_id, &u, &l, &p, &v]).unwrap();
        if let Ok(Some(row)) = rows.next() {
            return Some(json!({"eka": row.get::<&str, String>("eka").unwrap_or_default(), "dvi": row.get::<&str, String>("dvi").unwrap_or_default(), "bahu": row.get::<&str, String>("bahu").unwrap_or_default()}));
        }

        // Try dynamic sandhi generation
        if !u.is_empty() {
            let mut base_stmt = conn.prepare("SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND derivative='base' LIMIT 1").unwrap();
            let mut base_rows = base_stmt.query([&target_id, &l, &p, &v]).unwrap();
            if let Ok(Some(row)) = base_rows.next() {
                let eka = apply_upasarga_sandhi(&u, &row.get::<&str, String>("eka").unwrap_or_default());
                let dvi = apply_upasarga_sandhi(&u, &row.get::<&str, String>("dvi").unwrap_or_default());
                let bahu = apply_upasarga_sandhi(&u, &row.get::<&str, String>("bahu").unwrap_or_default());
                return Some(json!({"eka": eka, "dvi": dvi, "bahu": bahu, "note": "Dynamically Sandhi-fused"}));
            }
        }
        None
    }).await.unwrap();

    match row_json { Some(res) => res, None => json!({"error": "Verb combination not found in database."}) }
}

pub async fn generate_participle(pool: &Pool, dhatu_id: &str, root: &str, upasarga: &str, pratyaya: &str, gender: &str) -> Value {
    let d = dhatu_id.to_string(); let r = root.to_string(); let u = upasarga.to_string(); let pr = pratyaya.to_string(); let g = gender.to_string();
    
    // Clone `pr` before moving it into the closure so we can check it later
    let pr_clone = pr.clone();

    let base_form: Option<String> = pool.get().await.unwrap().interact(move |conn| {
        let target_id = if !d.is_empty() { d } else if !r.is_empty() {
            let mut stmt = conn.prepare("SELECT dhatu_id FROM info WHERE value = ?1 LIMIT 1").unwrap();
            let mut rows = stmt.query([&r]).unwrap();
            if let Ok(Some(row)) = rows.next() { row.get::<usize, String>(0).unwrap_or_default() } else { return None; }
        } else { return None; };

        // Try exact match
        let mut stmt = conn.prepare("SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga=?2 AND pratyaya=?3 AND derivative='base' LIMIT 1").unwrap();
        let mut rows = stmt.query([&target_id, &u, &pr_clone]).unwrap();
        if let Ok(Some(row)) = rows.next() {
            let form: String = row.get("base_form").unwrap_or_default();
            return Some(form.split(',').next().unwrap_or("").to_string());
        }

        // Try dynamic sandhi generation
        if !u.is_empty() {
            let mut base_stmt = conn.prepare("SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga='' AND pratyaya=?2 AND derivative='base' LIMIT 1").unwrap();
            let mut base_rows = base_stmt.query([&target_id, &pr_clone]).unwrap();
            if let Ok(Some(row)) = base_rows.next() {
                let bare_form: String = row.get("base_form").unwrap_or_default();
                let bare_first = bare_form.split(',').next().unwrap_or("").to_string();
                return Some(apply_upasarga_sandhi(&u, &bare_first));
            }
        }
        None
    }).await.unwrap();

    match base_form {
        Some(base) => {
            if ["tumun", "ktvA", "lyap", "Ramul"].contains(&pr.as_str()) { return json!({"base_form": base, "type": "avyaya"}); }
            match decline_noun(&base, &g) {
                Ok(declensions) => json!({"base_form": base, "declensions": declensions}),
                Err(e) => json!({"error": e.to_string()})
            }
        },
        None => json!({"error": "Participle combination not found."})
    }
}

pub fn generate_declension(base: &str, gender: &str) -> Value {
    match decline_noun(base, gender) {
        Ok(declensions) => json!({"base_form": base, "declensions": declensions}),
        Err(e) => json!({"error": e.to_string()})
    }
}
