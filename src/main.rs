use axum::{
    routing::get,
    Router,
    extract::{Path, State},
    response::Html,
    Json,
};
use deadpool_sqlite::{Config, Pool, Runtime};
use serde_json::Value;
use std::sync::Arc;
use std::env;

mod engine;

struct AppState {
    pool: Pool,
}

#[tokio::main]
async fn main() {
    println!("🚀 Booting Sanskrit Morphological Engine in Rust...");

    let cfg = Config::new("data/skt_morphology.db");
    let pool = cfg.create_pool(Runtime::Tokio1).unwrap();

    let app_state = Arc::new(AppState { pool });

    let app = Router::new()
        .route("/", get(serve_ui))
        .route("/api/analyze/:word", get(analyze_word))
        .with_state(app_state);

    // DYNAMIC PORT BINDING: Required for Render.com
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
