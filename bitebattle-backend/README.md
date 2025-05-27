# bitebattle Backend

bitebattle Backend is a Go-based backend application for the bitebattle project, providing RESTful APIs for user accounts, groups, polls, voting, notifications, and more.

---

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [Database & Migrations](#database--migrations)
- [API Endpoints](#api-endpoints)
- [Data Model (UML)](#data-model-uml)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/bitebattle-backend.git
   cd bitebattle-backend
   ```

2. **Install Go dependencies:**
   ```sh
   go mod tidy
   ```

3. **Copy and configure environment variables:**
   ```sh
   cp .env.example .env
   # Edit .env with your DB credentials
   ```

---

## Usage

### Local Development

- **Start Postgres with Docker Compose:**
  ```sh
  make up
  ```

- **Run database migrations:**
  ```sh
  make migrate
  ```

- **Run the backend server:**
  ```sh
  make run
  # or
  go run cmd/main.go
  ```

- **Full dev workflow:**
  ```sh
  make dev
  ```

- **Stop containers:**
  ```sh
  make stop
  ```

---

## Environment Variables

Create a `.env` file in the project root:

```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
```

These are used for database connection and migrations.

---

## Project Structure

```
bitebattle-backend
├── cmd/
│   └── main.go                # Application entry point
├── internal/
│   ├── account/               # Account logic (handlers, service)
│   ├── auth/                  # Auth middleware and logic
│   ├── group/                 # Group logic
│   ├── notification/          # Notification logic
│   ├── poll/                  # Poll logic
│   └── user/                  # User logic
├── migrations/                # SQL migration files
├── docker-compose.yml         # Postgres container config
├── Makefile                   # Dev commands
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # Project documentation
```

---

## Database & Migrations

- **Database:** PostgreSQL (managed via Docker Compose)
- **Migrations:** Use [golang-migrate](https://github.com/golang-migrate/migrate) with SQL files in `migrations/`.

**Example Postgres URL:**
```
postgresql://turan1393:securePassword123@localhost:5432/bitebattle_dev
```

**Migration commands:**
```sh
make migrate
```

**Migration file naming:**  
`000001_create_users_table.up.sql` / `000001_create_users_table.down.sql`

---

## API Endpoints

| Endpoint                  | Method | Description                    |
|---------------------------|--------|--------------------------------|
| `/api/account/`           | GET    | Get user profile               |
| `/api/account/groups`     | GET    | Get user groups                |
| `/api/account/polls`      | GET    | Get user polls                 |
| `/api/account/update`     | POST   | Update user profile            |
| `/api/auth/register`      | POST   | Register a new user            |
| `/api/auth/login`         | POST   | User login                     |
| `/api/groups`             | POST   | Create a group                 |
| `/api/groups/:id`         | GET    | Get group info                 |
| `/api/polls`              | POST   | Create a poll                  |
| `/api/polls/:id/vote`     | POST   | Cast a vote                    |
| `/api/notifications`      | GET    | Get user notifications         |
| ...                       | ...    | ...                            |

See `internal/*/handler.go` for full endpoint definitions.

---

## Data Model (UML)

### Users, Groups, Polls, Votes, Notifications

```
+---------------------+            +---------------------+
|       users         |            |      groups         |
+---------------------+            +---------------------+
| id (PK)             |◄────────┐  | id (PK)             |
| name                |         └──┤ created_by (FK→users.id)
| email (UNIQUE)      |            | name                |
| password_hash       |            | invite_code (UNIQUE)|
| created_at          |            | created_at          |
+---------------------+            +---------------------+

        ▲                                     ▲
        │                                     │
+---------------------+           +----------------------+
|   group_members     |           |      restaurants     |
+---------------------+           +----------------------+
| id (PK)             |           | id (PK)              |
| user_id (FK→users)  |────────┐  | yelp_id (UNIQUE)     |
| group_id (FK→groups)|──────┐│  | name                 |
| joined_at           |      ││  | photo_url            |
+---------------------+      ││  | menu_url             |
                             ││  | created_at           |
                             ││  +----------------------+
                             ││
                             ▼▼
                    +----------------------+
                    |       polls          |
                    +----------------------+
                    | id (PK)              |
                    | group_id (FK→groups) |
                    | created_by (FK→users)|
                    | is_active            |
                    | created_at           |
                    | updated_at           |
                    +----------------------+

                             ▲
                             │
                    +----------------------+
                    |     poll_votes       |
                    +----------------------+
                    | id (PK)              |
                    | poll_id (FK→polls)   |
                    | user_id (FK→users)   |
                    | option_id (FK→poll_options) |
                    | restaurant_place_id  |
                    | created_at           |
                    +----------------------+

+----------------------+
|   notifications      |
+----------------------+
| id (PK)              |
| user_id (FK→users)   |
| message              |
| read                 |
| created_at           |
+----------------------+
```

---

## Troubleshooting

- **Database connection errors:**  
  Ensure `.env` is correct and Postgres is running (`make up`).

- **Migration errors:**  
  - *Duplicate migration file*: Remove or rename duplicates in `migrations/`.
  - *Column does not exist*: Check your migration files and run `make migrate` again.

- **Type assertion errors (`interface {} is string, not uuid.UUID`):**  
  Always parse string IDs to `uuid.UUID` in handlers:
  ```go
  userIDStr := c.MustGet("userID").(string)
  userID, err := uuid.Parse(userIDStr)
  ```

- **Null value in NOT NULL column:**  
  Ensure all required fields are provided in your insert statements.

---

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any suggestions or improvements.

---

## License

This project is licensed under the MIT License. See the LICENSE file for more details.