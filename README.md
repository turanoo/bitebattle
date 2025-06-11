# BiteBattle

BiteBattle is a Go-based RESTful API server powering the BiteBattle application. It manages user authentication, polls, voting, notifications, head-to-head matches, and restaurant search, providing a robust backend for collaborative food decision-making.

---

## Installation and Local Development

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience)

### 1. Clone the Repository

```sh
git clone https://github.com/turanoo/bitebattle.git
cd bitebattle
```

### 2. Configure Environment Variables

Copy the example file and edit as needed:

```sh
cp .env.example .env
# Edit .env with your desired DB credentials and your Google Places API Key if available
```

**Example `.env`:**
```
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_HOST=localhost
DB_PORT=5432
JWT_SECRET=your_jwt_secret
GOOGLE_PLACES_API_KEY=your_google_places_api_key
```

### 3. Start PostgreSQL with Docker Compose

Ensure your Docker daemon is prior to executing the next steps!
Skip to step 6 if you want to run the next 3 commands in one go. 

```sh
make up
```

### 4. Run Database Migrations

```sh
make migrate
```

### 5. Run the Backend Server

```sh
make run
# or
go run cmd/server/main.go
```

### 6. Full Local Dev Workflow

```sh
make dev
```

### 7. Stopping and Cleaning Up

```sh
make stop      # Stop containers
make destroy   # Stop and remove volumes
```

---

## API Overview

See the code in `internal/*/handler.go` for full details.  
All available endpoints:

| Endpoint                                 | Method | Description                                 |
|-------------------------------------------|--------|---------------------------------------------|
| `/api/auth/register`                     | POST   | Register a new user                         |
| `/api/auth/login`                        | POST   | User login                                  |
| `/api/account`                           | GET    | Get user profile                            |
| `/api/account`                           | PUT    | Update user profile                         |
| `/api/users/:id`                         | GET    | Get user by ID                              |
| `/api/users/?email=`                     | GET    | Get user by email                           |
| `/api/polls`                             | POST   | Create a poll                               |
| `/api/polls`                             | GET    | List polls for user                         |
| `/api/polls/:pollId/join`                | POST   | Join a poll by invite code                  |
| `/api/polls/:pollId`                     | GET    | Get poll details                            |
| `/api/polls/:pollId`                     | DELETE | Delete a poll                               |
| `/api/polls/:pollId/options`             | POST   | Add one or more options to a poll           |
| `/api/polls/:pollId/vote`                | POST   | Cast a vote                                 |
| `/api/polls/:pollId/unvote`              | POST   | Remove a vote                               |
| `/api/polls/:pollId/results`             | GET    | Get poll results                            |
| `/api/notifications`                     | GET    | Get user notifications                      |
| `/api/h2h/match`                         | POST   | Create head-to-head match                   |
| `/api/h2h/match/:id/accept`              | POST   | Accept a head-to-head match                 |
| `/api/h2h/match/:id/swipe`               | POST   | Submit a swipe for a match                  |
| `/api/h2h/match/:id/results`             | GET    | Get mutual likes for a match                |
| `/api/restaurants/search`                | GET    | Search restaurants (Google Places)          |

---

## System Design

### Tech Stack

- **Language:** Go (1.24+)
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL (managed via Docker Compose)
- **ORM/DB Driver:** [lib/pq](https://github.com/lib/pq)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Authentication:** JWT (via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt))
- **Logging:** [logrus](https://github.com/sirupsen/logrus)
- **Restaurant Search:** Google Places API

### Project Structure

```
bitebattle/
├── api/                # Route setup
├── cmd/server/         # Main entrypoint
├── internal/
│   ├── account/        # User profile management
│   ├── auth/           # Auth logic and middleware
│   ├── head2head/      # Head-to-head match logic
│   ├── notification/   # Notification logic
│   ├── poll/           # Polls, options, votes
│   ├── restaurant/     # Restaurant search
│   └── user/           # User CRUD
├── migrations/         # SQL migration files
├── pkg/
│   ├── db/             # DB connection
│   ├── logger/         # Logging utilities
│   └── utils/          # Utility functions
├── scripts/            # Helper scripts (e.g., run_migrations.sh)
├── docker-compose.yml  # Postgres container config
├── Makefile            # Dev commands
├── go.mod
├── go.sum
└── README.md
```

### Key Features

- **JWT Authentication:** Secure endpoints with token-based auth.
- **Role-based Polls:** Poll creators are "owners", others are "members".
- **Head-to-Head Matches:** Invite and swipe for food matches.
- **Notifications:** In-app notification system.
- **Restaurant Search:** Google Places integration.
- **Robust Logging:** Centralized, sanitized logging with logrus.

---

## Database Schema UML

### Users, Polls, Poll Members, Poll Options, Poll Votes, Notifications

```
+---------------------+
|       users         |
+---------------------+
| id (PK, UUID)       |
| email (UNIQUE)      |
| name                |
| password_hash       |
| created_at (TS)     |
| updated_at (TS)     |
+---------------------+
         ▲
         │
+----------------------+
|   polls_members      |
+----------------------+
| id (PK, UUID)        |
| poll_id (FK→polls)   |
| user_id (FK→users)   |
| joined_at (TS)       |
| UNIQUE(poll_id, user_id)
+----------------------+
         ▲
         │
+----------------------+
|       polls          |
+----------------------+
| id (PK, UUID)        |
| name                 |
| invite_code (UNIQUE) |
| created_by (FK→users)|
| is_active (bool)     |
| created_at (TS)      |
| updated_at (TS)      |
+----------------------+
         │
         ▼
+----------------------+
|    poll_options      |
+----------------------+
| id (PK, UUID)        |
| poll_id (FK→polls)   |
| restaurant_id        |
| name                 |
| image_url            |
| menu_url             |
| UNIQUE(poll_id, restaurant_id)
+----------------------+
         │
         ▼
+----------------------+
|     poll_votes       |
+----------------------+
| id (PK, UUID)        |
| poll_id (FK→polls)   |
| option_id (FK→poll_options) |
| user_id (FK→users)   |
| created_at (TS)      |
| UNIQUE(poll_id, user_id, option_id)
+----------------------+

+----------------------+
|   notifications      |
+----------------------+
| id (PK, UUID)        |
| user_id (FK→users)   |
| message              |
| read (bool)          |
| created_at (TS)      |
+----------------------+
```

### Head2Head Matches & Swipes

```
+--------------------------+
|   head2head_matches      |
+--------------------------+
| id (PK, UUID)            |
| inviter_id (FK→users)    |
| invitee_id (FK→users)    |
| status                   |
| categories (TEXT[])      |
| created_at (TS)          |
| updated_at (TS)          |
+--------------------------+
        ▲
        │
+--------------------------+
|   head2head_swipes       |
+--------------------------+
| id (PK, UUID)            |
| match_id (FK→head2head_matches) |
| user_id (FK→users)       |
| restaurant_id            |
| restaurant_name          |
| liked (bool)             |
| created_at (TS)          |
+--------------------------+
```

---

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any suggestions or improvements.

---

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
