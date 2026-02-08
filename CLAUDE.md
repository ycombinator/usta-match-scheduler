# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Does

A USTA (United States Tennis Association) NorCal league match scheduling tool. It scrapes team/match data from the USTA NorCal website, lets users set team day preferences, generates optimized schedules respecting constraints, and provides a drag-and-drop calendar UI for manual adjustments.

## Development

Run the full stack (frontend + backend) in development mode with hot reload:
```bash
docker compose up --watch
```
Application available at `http://localhost:3000`. The web container proxies API requests to the backend on port 8000.

Set `USE_MOCK_DATA=true` env var to use mock data instead of live USTA scraping.

## Running Tests

Go tests (from repo root):
```bash
go test ./...                           # all tests
go test ./internal/scheduler/...        # single package
go test ./internal/models/... -run TestName  # single test
```
Tests use `testify` for assertions. USTA scraping tests use embedded HTML fixtures.

No frontend test framework is configured.

## Build Commands

```bash
# Backend (Go)
go build ./api/...

# Frontend
cd web && npm install && npx webpack --mode=production
```

## Architecture

### Backend (Go)

- **`api/`** — HTTP server on `:8000` using standard `net/http` ServeMux
  - `api/routing/router.go` — Three endpoints:
    - `GET /usta/organization/{id}/teams` — scrape teams from USTA site
    - `GET /usta/organization/{id}/matches` — scrape matches from USTA site
    - `POST /schedule` — generate schedule from team preferences
- **`internal/models/`** — Domain types: `Team`, `Event`, `Match`, `Schedule`, `Input`, `DaySlot`
- **`internal/scheduler/`** — Scheduling algorithms implementing a `Scheduler` interface with `Run()`. Three implementations: `NewPreferring()` (primary, respects day preferences), `NewEager()`, `NewConstraining()`. Schedules matches week-by-week respecting blackout dates, team isolation pairs, and daytime/evening constraints.
- **`internal/usta/`** — Web scraper using `goquery` to parse USTA NorCal league HTML pages for team rosters and match results
- **`main.go`** — CLI entry point (separate from the API server)
- **`input.yml`** — YAML config defining teams, blackout dates, scheduling types, and day preferences

### Frontend (React, in `web/`)

- **`web/src/App.jsx`** — Class-based root component. Two modes: "set_team_preferences" (configure teams) and "edit_schedule" (calendar view). Persists schedule to localStorage.
- **`web/src/components/`** — Calendar components (`CalendarMonthGroup` → `CalendarMonth` → `CalendarWeek` → `CalendarDay` → `CalendarEvent`), drag-and-drop via `@dnd-kit/core`, team preference UI (`TeamPreferences`, `OrderedSelectionGroup`)
- **`web/src/lib/date_utils.js`** — Date manipulation utilities (month/week calculations, event-to-day mapping)
- **`web/server.js`** — Express server that serves webpack output and proxies `/api/*` to the Go backend

### Data Flow

1. Frontend fetches teams from `GET /usta/organization/{id}/teams` (backend scrapes USTA site)
2. User configures day preferences and constraints in the UI
3. Frontend posts to `POST /schedule` with preferences
4. Backend runs scheduler algorithm, fetches match data from USTA, returns generated schedule
5. Frontend renders schedule in calendar with drag-and-drop editing
6. Schedule export via `html-to-image`

### Scheduling Constraints

- **Daytime teams**: weekday mornings only (no Saturday/Sunday)
- **Evening teams**: weekday evenings, flexible weekends
- **Blackout dates**: no matches scheduled (holidays, etc.)
- **Team isolation pairs**: teams that can't play in the same time slot
- **Day preferences**: ordered list of preferred match days per team
