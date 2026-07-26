# Stage 1: Build the Rust binary
FROM rust:1.75-slim-bookworm as builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 2: Create the lightweight production image
FROM debian:bookworm-slim
WORKDIR /app

# Install SQLite dependencies AND 'wget'
RUN apt-get update && apt-get install -y libsqlite3-dev ca-certificates wget && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /app/target/release/skt-morph-tool /app/skt-morph-tool

# Copy the HTML frontend
COPY index.html /app/index.html

# CREATE data folder and DOWNLOAD the DB directly from your GitHub Release!
RUN mkdir -p data && wget -q -O data/skt_morphology.db https://github.com/eadaradhiraj/skt-morph-tool/releases/download/v1.0/skt_morphology.db

EXPOSE 8000
CMD ["/app/skt-morph-tool"]