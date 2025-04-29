# Ingest Demo

[![Go](https://github.com/navicore/ingest-demo/actions/workflows/go.yml/badge.svg)](https://github.com/navicore/ingest-demo/actions/workflows/go.yml)

A CLI tool for generating and processing JSON data, with Parquet file output.

## Features

- Generate sample boat/vessel JSON data to stdout
- Process JSON data from stdin and save to partitioned Parquet files
- Data partitioning by name/year/month/day

## Installation

```bash
go install github.com/navicore/ingest-demo/cmd/ingest@latest
```

Or build from source:

```bash
git clone github.com/navicore/ingest-demo
cd ingest-demo
go build -o ingest ./cmd/ingest
```

## Usage

### Generate JSON data

Generate 100 sample JSON records:

```bash
ingest --generate --count 100
```

### Process JSON data to Parquet

Pipe generated data to the processor:

```bash
ingest --generate --count 100 | ingest
```

Or process data from a file:

```bash
cat data.json | ingest
```

### Specify output directory

```bash
ingest --generate --count 100 | ingest --output /path/to/output
```

## Output Format

The processor creates Parquet files organized in directories by:
- Boat/vessel name
- Year
- Month
- Day

Example path structure:
```
output/
  ├── Boat 1/
  │   └── 2025/
  │       └── 04/
  │           └── 28/
  │               └── data_a1b2c3d4.parquet
  ├── Boat 2/
  │   └── 2025/
  │       └── 04/
  │           └── 28/
  │               └── data_e5f6g7h8.parquet
```

## Data Format

### Input JSON Format

```json
{
  "version": "1.0.0",
  "name": "Boat 1",
  "uuid": "urn:mrn:signalk:uuid:182010c5-20f6-4c8a-9f32-fd18af029d3a",
  "latitude": 37.78039400669318,
  "longitude": -122.38526611923439,
  "altitude": 0.0,
  "course": 24.786300968644476,
  "speed": 6,
  "timestamp": "2025-04-16T07:18:45.592502Z"
}
```

### Output Parquet Schema

- version: string
- name: string
- uuid: string
- latitude: double
- longitude: double
- altitude: double
- course: double
- speed: double
- timestamp: string (ISO 8601)
- year: int32
- month: int32
- day: int32
- processed_at: string (ISO 8601)

## Development

### Requirements

- Go 1.18+
- Make (optional)

### Build and Test

Using Make:

```bash
# Build the project
make build

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean

# Run integration tests
make test-integration
```

Using Go directly:

```bash
# Build
go build -o ingest ./cmd/ingest

# Test
go test -v ./...

# Lint (requires golangci-lint)
golangci-lint run ./...
```

## GitHub Actions

This repository is configured with GitHub Actions for:

1. Building the project
2. Running tests
3. Linting the code

The workflow configuration is located in `.github/workflows/go.yml`.