# Ghana Location API

A read-only REST API providing normalized administrative location data for Ghana, including countries, regions, districts, constituencies, and cities.
https://ghana-location-api.vercel.app/api/v1/regions

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

Run the migration using the migrate command:

```bash
go run cmd/migrate/main.go
```

Or manually using psql:

```bash
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

### 5. Seed data

Seed the database with location data:

```bash
go run cmd/seed/main.go
```

This will load data from the JSON files in the `data/` directory:

- `countries.json`
- `regions.json`
- `districts.json`
- `constituencies.json`
- `cities.json`

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

- **Handlers** (`pkg/handlers/`) - HTTP layer only, no business logic
- **Services** (`pkg/services/`) - Business logic and validation
- **Repositories** (`pkg/repositories/`) - Database access using pgx
- **Models** (`pkg/models/`) - Domain entities

### Database Schema

- `countries` - Country information
- `regions` - Administrative regions
- `districts` - Districts, metros, and municipals
- `constituencies` - Electoral constituencies
- `cities` - Cities and towns with coordinates

All tables use UUID primary keys and slug fields for public identifiers. Slugs are stable and never change.

## Deployment

### Vercel (Serverless)

This API is configured for deployment on Vercel as a serverless function.

#### Prerequisites

- Vercel account
- PostgreSQL database (Supabase recommended)
- `DATABASE_URL environment variable

#### Deployment Steps

1. **Connect your repository to Vercel**

   - Import your GitHub repository in the Vercel dashboard
   - Vercel will auto-detect the Go configuration

2. **Set environment variables**

   - In your Vercel project settings, add:
     - `DATABASE_URL`: Your PostgreSQL connection string

3. **Run migrations and seed data**

   - Before deploying, ensure your production database is set up:

   ```bash
   # Set DATABASE_URL to your production database
   export DATABASE_URL="your-production-database-url"

   # Run migrations
   go run cmd/migrate/main.go

   # Seed data
   go run cmd/seed/main.go
   ```

4. **Deploy**
   - Push to your main branch or use `vercel deploy`
   - Vercel will automatically build and deploy your function

#### Configuration

The `vercel.json` file configures:

- Go runtime (auto-detected from `go.mod`)
- Routing to the serverless function
- Build settings

#### Important Notes

- The API uses `pkg/` instead of `internal/` packages to avoid Go's internal package visibility restrictions in serverless environments
- Go version is set to 1.24 (Vercel's current supported version)
- The handler is located at `api/index.go` and exports a `Handler` function

#### Live API

The deployed API is available at:

- Production: `https://ghana-location-api.vercel.app`
- Health check: `https://ghana-location-api.vercel.app/health`

## Technology Stack

- **Language**: Go 1.24
- **Router**: Chi
- **Database**: PostgreSQL with pgx driver
- **Configuration**: Environment variables via godotenv
- **Deployment**: Vercel (serverless)

## Development

### Project Structure

```
ghana-location-api/
├── api/
│   └── index.go            # Vercel serverless function handler
├── cmd/
│   ├── api/
│   │   └── main.go         # Local development server
│   ├── migrate/
│   │   └── main.go         # Database migration tool
│   └── seed/
│       └── main.go         # Database seeding tool
├── pkg/
│   ├── handlers/           # HTTP handlers
│   ├── services/           # Business logic
│   ├── repositories/       # Database access
│   ├── models/             # Domain models
│   ├── config/             # Configuration
│   └── errors/             # Error handling
├── migrations/             # SQL migrations
├── data/                   # Seed data files (JSON)
├── scripts/                # Data processing scripts
├── vercel.json             # Vercel deployment configuration
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
