# Umineko Quote Search

A quote search engine for Umineko no Naku Koro ni. Search through thousands of lines of dialogue from the visual novel.

## Features

- Fuzzy search through all dialogue
- Filter by character and episode
- Random quote generator
- English/Japanese language toggle
- Beautiful Umineko-themed web interface
- Single executable (all assets embedded)

## Quick Start

```bash
go build -o umineko_quote.exe .
./umineko_quote.exe
```

Open http://127.0.0.1:3000

## API Endpoints

| Endpoint                    | Description                                  |
|-----------------------------|----------------------------------------------|
| `GET /api/v1/search`        | Fuzzy search quotes                          |
| `GET /api/v1/random`        | Get random quote                             |
| `GET /api/v1/character/:id` | Get quotes by character ID                   |
| `GET /api/v1/characters`    | List all character IDs and names             |
| `GET /api/v1/health`        | Health check                                 |

### Query Parameters

| Parameter   | Endpoints                      | Description                              |
|-------------|--------------------------------|------------------------------------------|
| `q`         | search                         | Search query (required)                  |
| `lang`      | search, random, character      | Language: `en` (default) or `ja`         |
| `character` | search, random                 | Filter by character ID                   |
| `episode`   | search, random, character      | Filter by episode (1-8)                  |
| `limit`     | search, character              | Results per page (default: 30)           |
| `offset`    | search, character              | Pagination offset                        |

### Response Format

```json
{
  "results": [
    {
      "quote": {
        "text": "Without love, it cannot be seen.",
        "characterId": "27",
        "character": "Beatrice",
        "audioId": "10700001",
        "episode": 1
      },
      "score": 95
    }
  ]
}
```

## Build

### Windows
```powershell
go build -o umineko_quote.exe .
```

### Linux
```bash
go build -o umineko_quote .
```

### Cross-compile

```powershell
# Mac ARM (M1/M2/M3)
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o umineko_quote_mac .; $env:GOOS=""; $env:GOARCH=""

# Mac Intel
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o umineko_quote_mac_intel .; $env:GOOS=""; $env:GOARCH=""

# Linux x64
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o umineko_quote_linux .; $env:GOOS=""; $env:GOARCH=""
```

## Docker

```bash
docker build -t umineko-quote .
docker run -p 3000:3000 umineko-quote
```

## Data

Quote data is parsed from Umineko no Naku Koro ni script files. Place language files in `internal/quote/data/` before building:

```
internal/quote/data/
├── english.txt
└── japanese.txt
```

Additional languages can be added by placing new files in this directory and updating the `langFiles` map in `service.go`.
