use axum::{
    routing::get,
    Router,
    extract::{Path, State, Query},
    response::Html,
    Json,
};
use deadpool_sqlite::{Config, Pool, Runtime};
use serde_json::Value;
use serde::Deserialize;
use std::sync::Arc;
use std::env;

mod engine;

struct AppState {
    pool: Pool,
}

#[derive(Deserialize)]
struct ConjugateQuery {
    #[serde(default)] dhatu_id: String,
    #[serde(default)] root: String,
    #[serde(default)] upasarga: String,
    #[serde(default = "default_lakara")] lakara: String,
    #[serde(default = "default_purusha")] purusha: String,
    #[serde(default = "default_voice")] voice: String,
}

fn default_lakara() -> String { "law".to_string() }
fn default_purusha() -> String { "praTama".to_string() }
fn default_voice() -> String { "parasmEpadam".to_string() }

#[tokio::main]
async fn main() {
    println!("🚀 Booting Sanskrit Morphological Engine in Rust...");

    let cfg = Config::new("data/skt_morphology.db");
    let pool = cfg.create_pool(Runtime::Tokio1).unwrap();

    let app_state = Arc::new(AppState { pool });

    let app = Router::new()
        .route("/", get(serve_ui))
        .route("/api/analyze/:word", get(analyze_word))
        .route("/api/dhatus/:query", get(search_dhatu_api))
        .route("/api/generate/verb", get(generate_verb_api))
        .with_state(app_state);

    let port = env::var("PORT").unwrap_or_else(|_| "8000".to_string());
    let addr = format!("0.0.0.0:{}", port);

    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    println!("✅ Server running on http://{}", addr);
    axum::serve(listener, app).await.unwrap();
}

async fn serve_ui() -> Html<&'static str> {
    Html(include_str!("../index.html"))
}

async fn analyze_word(
    Path(word): Path<String>,
    State(state): State<Arc<AppState>>,
) -> Json<Value> {
    let results = engine::analyzer::analyze(&word, &state.pool).await;
    Json(results)
}

async fn search_dhatu_api(
    Path(query): Path<String>,
    State(state): State<Arc<AppState>>,
) -> Json<Value> {
    let results = engine::dhatu::search_dhatu(&query, &state.pool).await;
    Json(results)
}

async fn generate_verb_api(
    Query(params): Query<ConjugateQuery>,
    State(state): State<Arc<AppState>>,
) -> Json<Value> {
    let results = engine::generator::generate_verb(
        &state.pool, 
        &params.dhatu_id, 
        &params.root,
        &params.upasarga, 
        &params.lakara, 
        &params.purusha, 
        &params.voice
    ).await;
    Json(results)
}
