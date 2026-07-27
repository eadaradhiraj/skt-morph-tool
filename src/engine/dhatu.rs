use deadpool_sqlite::Pool;
use serde_json::{json, Value};
use std::collections::HashMap;

pub async fn search_dhatu(query: &str, pool: &Pool) -> Value {
    let q_clone = format!("%{}%", query);
    
    let results: Vec<Value> = pool.get().await.unwrap().interact(move |conn| {
        // Find all dhatu_ids where any metadata (name, meaning) matches the query
        let mut stmt = conn.prepare(
            "SELECT dhatu_id, key_name, value FROM info 
             WHERE dhatu_id IN (SELECT dhatu_id FROM info WHERE value LIKE ?1 LIMIT 100)"
        ).unwrap();
        
        let mut rows = stmt.query([&q_clone]).unwrap();
        let mut map: HashMap<String, serde_json::Map<String, Value>> = HashMap::new();
        
        while let Ok(Some(row)) = rows.next() {
            let id: String = row.get(0).unwrap_or_default();
            let key: String = row.get(1).unwrap_or_default();
            let val: String = row.get(2).unwrap_or_default();
            
            let entry = map.entry(id.clone()).or_insert_with(serde_json::Map::new);
            entry.insert(key, json!(val));
            entry.insert("dhatu_id".to_string(), json!(id));
        }
        
        map.into_values().map(Value::Object).collect()
    }).await.unwrap();

    json!(results)
}
