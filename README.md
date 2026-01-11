# Ghana Location API

A read-only REST API providing normalized administrative location data for Ghana, including countries, regions, districts, constituencies, and cities.

## Overview

This API serves as a canonical source for Ghana's administrative geography hierarchy:

```
Country
 └── Regions
      └── Districts / Metros / Municipals
           └── Constituencies
                └── Cities / Towns
```

The API is designed to be:
- **Read-only** - No mutations, no user data
- **Fast** - Optimized queries with proper indexing
- **Stable** - Slug-based identifiers that never change
- **Cache-friendly** - HTTP cache headers on all responses

## Prerequisites

- Go 1.22 or higher
- PostgreSQL (or Supabase)
- Supabase CLI (if using Supabase)

## Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/blingyplus/ghana-location-api.git
cd ghana-location-api
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgres://user:password@localhost:5432/ghana_location
PORT=8080
```

For Supabase, you can get your connection string from your Supabase project dashboard or use:

```bash
supabase db get-connection-string
```

### 4. Run database migrations

Connect to your database and run the migration:

```bash
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

Or using Supabase CLI:

```bash
supabase db reset
# Then manually run the migration file
```

### 5. Seed data (optional)

Place normalized data files in the `data/` directory:
- `regions.json`
- `districts.json`
- `constituencies.json`
- `cities.json`

A seeder script can be created to load this data into the database.

### 6. Run the API

```bash
go run cmd/api/main.go
```

The API will start on port 8080 (or the port specified in your `.env` file).

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Countries

- `GET /api/v1/countries` - List all countries
- `GET /api/v1/countries/{code}` - Get country by code (e.g., "GH")

### Regions

- `GET /api/v1/regions` - List all regions
- `GET /api/v1/regions/{slug}` - Get region by slug
- `GET /api/v1/regions/{slug}/districts` - Get districts in a region

### Districts

- `GET /api/v1/districts/{slug}` - Get district by slug
- `GET /api/v1/districts/{slug}/constituencies` - Get constituencies in a district

### Constituencies

- `GET /api/v1/constituencies/{slug}` - Get constituency by slug

### Cities

- `GET /api/v1/cities?district={slug}` - Get cities in a district

## Response Format

All responses are JSON. Success responses include cache headers:

```
Cache-Control: public, max-age=3600
Content-Type: application/json
```

### Example Response

```json
{
  "id": "uuid",
  "name": "Greater Accra",
  "slug": "greater-accra",
  "capital": "Accra"
}
```

### Error Response

```json
{
  "error": "resource not found"
}
```

## Architecture

The API follows clean architecture principles with clear boundaries:

```
Request → Handler → Service → Repository → Database
         ↓
      JSON Response
```

### Layers

- **Handlers** (`internal/handlers/`) - HTTP layer only, no business logic
- **Services** (`internal/services/`) - Business logic and validation
- **Repositories** (`internal/repositories/`) - Database access using pgx
- **Models** (`internal/models/`) - Domain entities

### Database Schema

- `countries` - Country information
- `regions` - Administrative regions
- `districts` - Districts, metros, and municipals
- `constituencies` - Electoral constituencies
- `cities` - Cities and towns with coordinates

All tables use UUID primary keys and slug fields for public identifiers. Slugs are stable and never change.

## Technology Stack

- **Language**: Go 1.22+
- **Router**: Chi
- **Database**: PostgreSQL with pgx driver
- **Configuration**: Environment variables via godotenv

## Development

### Project Structure

```
ghana-location-api/
├── cmd/api/
│   └── main.go              # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── services/            # Business logic
│   ├── repositories/       # Database access
│   ├── models/              # Domain models
│   ├── config/              # Configuration
│   └── errors/              # Error handling
├── migrations/              # SQL migrations
├── data/                    # Seed data files
├── scripts/                 # Data processing scripts
├── go.mod
└── README.md
```

### Building

```bash
go build -o bin/api cmd/api/main.go
```

### Running Tests

```bash
go test ./...
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

This is a read-only API designed for stability. Data normalization and seeding scripts are welcome contributions.
