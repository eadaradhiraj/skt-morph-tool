use axum::{routing::get, Router, extract::{Path, State, Query}, Json};
use tower_http::services::ServeDir;
use deadpool_sqlite::{Config, Pool, Runtime};
use serde_json::Value;
use serde::Deserialize;
use std::sync::Arc;
use std::env;

mod engine;
struct AppState { pool: Pool }

#[derive(Deserialize)]
struct AnalyzeQuery { #[serde(default = "d_page")] page: i64 }
#[derive(Deserialize)]
struct VerbQuery { root: String, #[serde(default)] upasarga: String, lakara: String, purusha: String, voice: String, prayoga: String, derivative: String }
#[derive(Deserialize)]
struct ParticipleQuery { root: String, #[serde(default)] upasarga: String, pratyaya: String, gender: String, derivative: String }
#[derive(Deserialize)]
struct DeclensionQuery { base: String, gender: String }

fn d_page() -> i64 { 1 }

#[tokio::main]
async fn main() {
    let cfg = Config::new("data/skt_morphology.db");
    let pool = cfg.create_pool(Runtime::Tokio1).unwrap();
    let app = Router::new()
        .nest_service("/", ServeDir::new("frontend/dist"))
        .route("/api/analyze/:word", get(analyze_word))
        .route("/api/dhatus/:query", get(search_dhatu_api))
        .route("/api/generate/verb", get(gen_verb))
        .route("/api/generate/participle", get(gen_participle))
        .route("/api/generate/declension", get(gen_declension))
        .with_state(Arc::new(AppState { pool }));
    
    let port = env::var("PORT").unwrap_or_else(|_| "8000".to_string());
    let addr = format!("0.0.0.0:{}", port);
    println!("✅ Server running on http://{}", addr);
    axum::serve(tokio::net::TcpListener::bind(&addr).await.unwrap(), app).await.unwrap();
}

async fn analyze_word(Path(word): Path<String>, Query(q): Query<AnalyzeQuery>, State(state): State<Arc<AppState>>) -> Json<Value> { Json(engine::analyzer::analyze(&word, q.page, &state.pool).await) }
async fn search_dhatu_api(Path(query): Path<String>, State(state): State<Arc<AppState>>) -> Json<Value> { Json(json!(engine::dhatu::search_dhatu(&query, &state.pool).await)) }
async fn gen_verb(Query(p): Query<VerbQuery>, State(s): State<Arc<AppState>>) -> Json<Value> { Json(engine::generator::generate_verb(&s.pool, &p.root, &p.upasarga, &p.lakara, &p.purusha, &p.voice, &p.prayoga, &p.derivative).await) }
async fn gen_participle(Query(p): Query<ParticipleQuery>, State(s): State<Arc<AppState>>) -> Json<Value> { Json(engine::generator::generate_participle(&s.pool, &p.root, &p.upasarga, &p.pratyaya, &p.gender, &p.derivative).await) }
async fn gen_declension(Query(p): Query<DeclensionQuery>) -> Json<Value> { Json(engine::generator::generate_declension(&p.base, &p.gender)) }
