use deadpool_sqlite::Pool;
use serde_json::{json, Value};

pub async fn search_dhatu(query: &str, pool: &Pool) -> Vec<Value> {
    let q = query.to_string();
    pool.get().await.unwrap().interact(move |conn| {
        let mut res = Vec::new();
        let search_term = format!("%{}%", q);
        
        // This query finds the Root, Meaning, and generated forms (like 'Adara') just like the screenshot!
        let mut stmt = conn.prepare("
            SELECT p.dhatu_id, p.pratyaya, p.base_form, p.upasarga,
                   (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('mUlaDAtuH', 'DAtuH') LIMIT 1) as root,
                   (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('aTfaH', 'meaning', 'eng', 'hin') LIMIT 1) as meaning,
                   (SELECT value FROM info WHERE dhatu_id = p.dhatu_id AND key_name IN ('gaRaH', 'gana') LIMIT 1) as gana
            FROM participles p
            WHERE p.base_form LIKE ?1 OR root LIKE ?1 OR meaning LIKE ?1
            LIMIT 100
        ").unwrap();
        
        let mut rows = stmt.query([&search_term]).unwrap();
        while let Ok(Some(row)) = rows.next() {
            res.push(json!({
                "dhatu_id": row.get::<&str, String>("dhatu_id").unwrap_or_default(),
                "root": row.get::<&str, String>("root").unwrap_or_default(),
                "meaning": row.get::<&str, String>("meaning").unwrap_or_default(),
                "gana": row.get::<&str, String>("gana").unwrap_or_default(),
                "pratyaya": row.get::<&str, String>("pratyaya").unwrap_or_default(),
                "upasarga": row.get::<&str, String>("upasarga").unwrap_or_default(),
                "base_form": row.get::<&str, String>("base_form").unwrap_or_default()
            }));
        }
        res
    }).await.unwrap()
}
