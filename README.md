# Light House

A lightweight uptime monitoring tool for tracking the availability and performance of web services and APIs.

## Features

- **Uptime Checks** — Monitor HTTP/HTTPS endpoints with configurable intervals
- **Background Worker** — Automated check execution with result tracking
- **Dashboard** — Clean UI showing check status, response times, and history
- **Multi-Tenant** — Organization-scoped data isolation
- **Auth** — JWT-based authentication with secure password hashing

## Quick Start

```bash
# Clone and start
git clone https://github.com/oFuterman/Light-House.git
cd Light-House
docker-compose up -d

# Access
# Frontend: http://localhost:3000
# API: http://localhost:8080/api/v1
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22, Fiber, GORM |
| Frontend | Next.js 14, TypeScript, Tailwind CSS |
| Database | PostgreSQL 16 |
| Infrastructure | Docker Compose |

## Documentation

- [MVP Specification](docs/mvp-spec.md) — Current features and technical details
- [Demo Walkthrough](docs/demo-walkthrough.md) — Guided demo script for presentations

## License

Apache 2.0 — See [LICENSE](LICENSE/LICENSE-2.0.txt)
