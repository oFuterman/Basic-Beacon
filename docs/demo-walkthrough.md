# Light House - Demo Walkthrough

## What is Light House?

Light House is a lightweight uptime monitoring tool that tracks the availability and performance of web services. It features automated background checks, historical result tracking, and a clean dashboard for managing monitored endpoints.

---

## Running Locally

### Prerequisites
- Docker and Docker Compose installed
- Ports 3000, 8080, and 5432 available

### Quick Start

```bash
# Clone the repository
git clone https://github.com/oFuterman/Light-House.git
cd Light-House

# Start all services
docker-compose up -d

# Wait ~30 seconds for services to initialize
```

**Access the app:**
- Frontend: http://localhost:3000
- API: http://localhost:8080/api/v1

---

## Demo Flow

### 1. Sign Up

1. Navigate to http://localhost:3000
2. Click **Sign up** (or go to `/signup`)
3. Enter:
   - Organization name: `Demo Company`
   - Email: `demo@example.com`
   - Password: `password123`
4. Click **Sign Up**

You'll be redirected to the dashboard.

### 2. Empty State

On the dashboard, observe the empty state:
- "No checks yet" message
- Explanation of what Light House does
- Primary CTA to create your first check

*This demonstrates thoughtful UX for new users.*

### 3. Create a Passing Check

1. Click **Create Your First Check** (or **New Check**)
2. Fill in:
   - Name: `Example.com`
   - URL: `https://example.com`
   - Interval: `60` seconds
3. Click **Create Check**

You'll be redirected to the check detail page.

### 4. Create a Failing Check (Optional)

1. Click **New Check** from the dashboard
2. Fill in:
   - Name: `Intentional 500`
   - URL: `https://httpbin.org/status/500`
   - Interval: `60` seconds
3. Click **Create Check**

*This demonstrates how Light House handles failing endpoints.*

### 5. Wait for Results

The background worker runs every 30 seconds. After 30–60 seconds:

1. Refresh the dashboard
2. Observe:
   - Status badges showing UP (green) or DOWN (red)
   - Last check timestamps
   - Next check times

3. Click on a check to view the detail page:
   - Check configuration (URL, interval, status)
   - Results table with response times and status codes
   - Error messages for failed checks

### 6. Logs Page (Under Construction)

1. Click **Logs** in the navigation
2. Observe:
   - "Under construction" banner
   - Preview of what the log viewer will look like
   - Placeholder log entries showing the planned format

*Explain: "Log ingestion via API keys is planned for Phase 9–10, along with alerting and notifications."*

### 7. Logout and Session Verification

1. Click **Logout**
2. You're redirected to the login page
3. Log back in — note the "Verifying session..." spinner
4. Refresh any protected page — the app verifies your token with the backend

*This demonstrates proper auth handling, not just trusting localStorage.*

---

## Talking Points

Use these to guide the conversation with interviewers:

### 1. Authentication & Security
- JWT-based auth with bcrypt password hashing
- Token verification on page refresh (calls `GET /me`)
- Protected routes with auth guards

### 2. Multi-Tenancy
- Organization-scoped data isolation
- Users only see their own checks and results
- Foundation for team features in future phases

### 3. Background Worker
- Go goroutine runs every 30 seconds
- Finds due checks based on `last_checked_at + interval`
- Stores results with status code, response time, and errors
- Updates check status automatically

### 4. Clean API Design
- RESTful endpoints following standard conventions
- Consistent error handling and response formats
- JWT middleware for protected routes

### 5. Frontend Architecture
- Next.js 14 with App Router
- TypeScript throughout
- Reusable components (Loading, ErrorState, StatusBadge)
- Context-based auth state management

### 6. Production Considerations
- Docker Compose for easy local development
- Environment-based configuration
- Graceful error handling and loading states
- Roadmap for tests, CI/CD, and deployment (Phase 8–10)

---

## Roadmap Preview

If asked about future plans:

| Phase | Focus |
|-------|-------|
| 8–10 | Tests, CI/CD, hosted demo |
| 11–14 | Stripe billing, team management, notifications |
| 15–18 | AI-powered check creation, repo scanning, monitoring advisor |

See [mvp-spec.md](./mvp-spec.md) for detailed specifications.
