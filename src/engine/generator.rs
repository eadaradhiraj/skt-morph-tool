use deadpool_sqlite::Pool;
use serde_json::{json, Value};

pub async fn generate_verb(
    pool: &Pool,
    dhatu_id: &str,
    root: &str,
    upasarga: &str,
    lakara: &str,
    purusha: &str,
    voice: &str,
) -> Value {
    let d = dhatu_id.to_string();
    let r = root.to_string();
    let u = upasarga.to_string();
    let l = lakara.to_string();
    let p = purusha.to_string();
    let v = voice.to_string();

    let row_json: Option<Value> = pool.get().await.unwrap().interact(move |conn| {
        
        // 1. Resolve ID if the user provided a text 'root' instead of 'dhatu_id'
        let target_id = if !d.is_empty() {
            d
        } else if !r.is_empty() {
            let mut stmt = conn.prepare("SELECT dhatu_id FROM info WHERE value = ?1 LIMIT 1").unwrap();
            let mut rows = stmt.query([&r]).unwrap();
            if let Ok(Some(row)) = rows.next() {
                row.get::<usize, String>(0).unwrap_or_default()
            } else {
                return Some(json!({"error": format!("Root '{}' not found in database.", r)}));
            }
        } else {
            return Some(json!({"error": "You must provide either 'dhatu_id' or 'root'."}));
        };

        // 2. Fetch the conjugation
        let mut stmt = conn.prepare(
            "SELECT eka, dvi, bahu FROM conjugations 
             WHERE dhatu_id=?1 AND upasarga=?2 AND lakara=?3 AND purusha=?4 AND voice=?5 AND derivative='base' LIMIT 1"
        ).unwrap();
        
        let mut rows = stmt.query([&target_id, &u, &l, &p, &v]).unwrap();
        
        if let Ok(Some(row)) = rows.next() {
            Some(json!({
                "dhatu_id_used": target_id,
                "eka": row.get::<&str, String>("eka").unwrap_or_default(),
                "dvi": row.get::<&str, String>("dvi").unwrap_or_default(),
                "bahu": row.get::<&str, String>("bahu").unwrap_or_default(),
            }))
        } else {
            None
        }
    }).await.unwrap();

    match row_json {
        Some(res) => res,
        None => json!({"error": "Verb combination not found in database. Check Voice (parasmEpadam/Atmanepadam) or Lakara."})
    }
}
