# Stage 1: Build the Rust binary
FROM rust:1.75-slim-bookworm as builder
WORKDIR /app
COPY . .
# Compile the app for release (optimized for max speed)
RUN cargo build --release

# Stage 2: Create the lightweight production image
FROM debian:bookworm-slim
WORKDIR /app

# Install required SQLite dependencies
RUN apt-get update && apt-get install -y libsqlite3-dev ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /app/target/release/skt-morph-tool /app/skt-morph-tool

# Copy the massive database and the HTML frontend
COPY data/skt_morphology.db /app/data/skt_morphology.db
COPY index.html /app/index.html

# CREATE data folder and DOWNLOAD the 543MB database directly from GitHub Releases!
RUN mkdir -p data && wget -q -O data/skt_morphology.db https://github.com/eadaradhiraj/skt-morph-tool/releases/download/v1.0/skt_morphology.db

# Expose port 8000 for the cloud provider
EXPOSE 8000

# Start the engine
CMD ["/app/skt-morph-tool"]
