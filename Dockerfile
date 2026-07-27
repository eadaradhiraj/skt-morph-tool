# Stage 1: Build the Frontend (Vite)
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Rust binary
FROM rust:latest as rust-builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 3: Create the lightweight production image
FROM debian:bookworm-slim
WORKDIR /app

# Install SQLite dependencies AND 'wget' to download our database
RUN apt-get update && apt-get install -y libsqlite3-dev ca-certificates wget && rm -rf /var/lib/apt/lists/*

# Copy the compiled Rust binary
COPY --from=rust-builder /app/target/release/skt-morph-tool /app/skt-morph-tool

# Copy the compiled Frontend (Vite dist folder)
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# CREATE data folder and DOWNLOAD the 543MB database directly from GitHub Releases!
RUN mkdir -p data && wget -q -O data/skt_morphology.db "https://github.com/eadaradhiraj/skt-morph-tool/releases/download/v1.0/skt_morphology.db"

EXPOSE 8000
CMD ["/app/skt-morph-tool"]
