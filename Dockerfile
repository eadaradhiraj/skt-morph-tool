# Stage 1: Build the Frontend (Vite)
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go binary
FROM golang:alpine AS go-builder
WORKDIR /app

# Copy module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source code & engine package
COPY . .

# Build a static Go binary (CGO_ENABLED=0 because modernc.org/sqlite is pure Go)
RUN CGO_ENABLED=0 GOOS=linux go build -o skt-morph-tool main.go

# Stage 3: Create the lightweight production image
FROM alpine:latest
WORKDIR /app
ENV GOMEMLIMIT=400MiB
ENV GOGC=100

# Install CA certificates, wget, python3 and sqlite to build the database from skt-morph-data
RUN apk add --no-cache ca-certificates wget python3 sqlite

# Copy the compiled Go binary from Stage 2
COPY --from=go-builder /app/skt-morph-tool /app/skt-morph-tool

# Copy the compiled Frontend (Vite dist folder) from Stage 1
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# DOWNLOAD skt-morph-data source archive and BUILD the SQLite database
RUN mkdir -p data && \
    wget -q -O /tmp/skt-morph-data.tar.gz "https://github.com/eadaradhiraj/skt-morph-data/archive/refs/tags/v0.1.0.tar.gz" && \
    tar -xzf /tmp/skt-morph-data.tar.gz -C /tmp && \
    python3 /tmp/skt-morph-data-0.1.0/scripts/json_to_sqlite.py --data-dir /tmp/skt-morph-data-0.1.0/data --output data/skt_morphology.db --force && \
    sqlite3 data/skt_morphology.db "CREATE INDEX IF NOT EXISTS idx_conj_form_value ON conjugation_forms(form_value); CREATE INDEX IF NOT EXISTS idx_part_base ON participle_forms(base); CREATE INDEX IF NOT EXISTS idx_part_m ON participle_forms(m); CREATE INDEX IF NOT EXISTS idx_part_f ON participle_forms(f); CREATE INDEX IF NOT EXISTS idx_part_n ON participle_forms(n); VACUUM;" && \
    rm -rf /tmp/skt-morph-data.tar.gz /tmp/skt-morph-data-0.1.0

EXPOSE 8000
CMD ["/app/skt-morph-tool"]