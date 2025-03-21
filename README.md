# DataILama

DataILama is a tool for managing and analyzing research papers using AI. It provides a web interface and API for searching, analyzing, and managing academic papers.

## Features

- Web interface for paper management
- AI-powered paper analysis
- REST API for programmatic access
- Vector search capabilities

## Installation

### From Release

1. Download the latest release from the [releases page](https://github.com/hydrocode-de/datailama/releases)
2. Choose the appropriate binary for your system
3. Make it executable (Linux/macOS):
   ```bash
   chmod +x datailama-<os>-<arch>
   ```
4. Run the server:
   ```bash
   ./datailama-<os>-<arch> serve
   ```

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/hydrocode-de/datailama.git
   cd datailama
   ```

2. Build the frontend:
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

3. Build the Go binary:
   ```bash
   go build
   ```

4. Run the server:
   ```bash
   ./datailama serve
   ```

## Configuration

The server can be configured using environment variables or command-line flags:

- `DATAILAMA_PORT`: Port to listen on (default: 8080)
- `DATAILAMA_DB_URL`: Database connection URL
- `DATAILAMA_OLLAMA_URL`: Ollama server URL

## Development

### Prerequisites

- Go 1.23 or later
- Node.js 20 or later
- PostgreSQL database
- Ollama server

### Running Tests

```bash
# Build frontend
cd frontend && npm install && npm run build && cd ..

# Run Go tests
go test -v ./...
```
