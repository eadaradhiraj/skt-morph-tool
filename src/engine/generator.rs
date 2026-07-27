use deadpool_sqlite::Pool;
use serde_json::{json, Value};
use crate::engine::sandhi::apply_upasarga_sandhi;
use crate::engine::declension::decline_noun;

pub async fn generate_verb(pool: &Pool, root: &str, upasarga: &str, lakara: &str, purusha: &str, voice: &str, prayoga: &str, derivative: &str) -> Value {
    let r = root.to_string(); let u = upasarga.to_string(); let l = lakara.to_string(); let p = purusha.to_string(); let v = voice.to_string(); let pr = prayoga.to_string(); let dev = derivative.to_string();

    let row_json: Option<Value> = pool.get().await.unwrap().interact(move |conn| {
        let is_id = r.chars().next().map_or(false, |c| c.is_ascii_digit());
        
        // Exact match
        let q_exact = if is_id {
            "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
        } else {
            "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND prayoga=?6 AND derivative=?7 LIMIT 1"
        };
        let mut stmt = conn.prepare(q_exact).unwrap();
        let mut rows = stmt.query([&r, &u, &l, &p, &v, &pr, &dev]).unwrap();
        if let Ok(Some(row)) = rows.next() {
            return Some(json!({"eka": row.get::<&str, String>("eka").unwrap_or_default(), "dvi": row.get::<&str, String>("dvi").unwrap_or_default(), "bahu": row.get::<&str, String>("bahu").unwrap_or_default()}));
        }

        // Dynamic sandhi generation
        if !u.is_empty() {
            let q_dyn = if is_id {
                "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id=?1 AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
            } else {
                "SELECT eka, dvi, bahu FROM conjugations WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga='' AND lakara=?2 AND purusha=?3 AND voice=?4 AND prayoga=?5 AND derivative=?6 LIMIT 1"
            };
            let mut base_stmt = conn.prepare(q_dyn).unwrap();
            let mut base_rows = base_stmt.query([&r, &l, &p, &v, &pr, &dev]).unwrap();
            if let Ok(Some(row)) = base_rows.next() {
                let eka = apply_upasarga_sandhi(&u, &row.get::<&str, String>("eka").unwrap_or_default());
                let dvi = apply_upasarga_sandhi(&u, &row.get::<&str, String>("dvi").unwrap_or_default());
                let bahu = apply_upasarga_sandhi(&u, &row.get::<&str, String>("bahu").unwrap_or_default());
                return Some(json!({"eka": eka, "dvi": dvi, "bahu": bahu, "note": "Dynamically Sandhi-fused"}));
            }
        }
        None
    }).await.unwrap();

    match row_json { Some(res) => res, None => json!({"error": "Verb combination not found. Ensure root, Voice, and Prayoga are compatible."}) }
}

pub async fn generate_participle(pool: &Pool, root: &str, upasarga: &str, pratyaya: &str, gender: &str, derivative: &str) -> Value {
    let r = root.to_string(); let u = upasarga.to_string(); let pr = pratyaya.to_string(); let g = gender.to_string(); let dev = derivative.to_string();
    let pr_clone = pr.clone();

    let base_form: Option<String> = pool.get().await.unwrap().interact(move |conn| {
        let is_id = r.chars().next().map_or(false, |c| c.is_ascii_digit());

        let q_exact = if is_id {
            "SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
        } else {
            "SELECT base_form FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga=?2 AND pratyaya=?3 AND derivative=?4 LIMIT 1"
        };
        let mut stmt = conn.prepare(q_exact).unwrap();
        let mut rows = stmt.query([&r, &u, &pr_clone, &dev]).unwrap();
        if let Ok(Some(row)) = rows.next() {
            return Some(row.get::<&str, String>("base_form").unwrap_or_default().split(',').next().unwrap_or("").to_string());
        }

        if !u.is_empty() {
            let q_dyn = if is_id {
                "SELECT base_form FROM participles WHERE dhatu_id=?1 AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
            } else {
                "SELECT base_form FROM participles WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value=?1) AND upasarga='' AND pratyaya=?2 AND derivative=?3 LIMIT 1"
            };
            let mut base_stmt = conn.prepare(q_dyn).unwrap();
            let mut base_rows = base_stmt.query([&r, &pr_clone, &dev]).unwrap();
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
