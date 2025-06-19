# BiteBattle

BiteBattle is a Go-based RESTful API server powering the BiteBattle application. It manages user authentication, polls, voting, notifications, head-to-head matches, and restaurant search, providing a robust backend for collaborative food decision-making.

---

## API Overview

### System Design

#### Tech Stack

- **Language:** Go (1.24+)
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL (managed via Docker Compose)
- **ORM/DB Driver:** [lib/pq](https://github.com/lib/pq)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Authentication:** JWT (via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt))
- **Logging:** [logrus](https://github.com/sirupsen/logrus)
- **Restaurant Search:** Google Places API

#### Key Features

- **JWT Authentication:** Secure endpoints with token-based auth.
- **Role-based Polls:** Poll creators are "owners", others are "members".
- **Head-to-Head Matches:** Invite and swipe for food matches.
- **Notifications:** In-app notification system.
- **Restaurant Search:** Google Places integration.
- **Robust Logging:** Centralized, sanitized logging with logrus.

---

### Interactive API Docs

You can preview the OpenAPI (Swagger) specification using this [link](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/turanoo/bitebattle/master/docs/api-spec.yaml).

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

Create an env file and add the following information:

```sh
touch .env
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
GCS_PROFILE_BUCKET=your_gcs_profile_pictures_storage_bucket
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


## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any suggestions or improvements.

### Development Guidelines

Before opening a merge request (MR) against the master branch, ensure that the following commands complete successfully:

```sh
make lint   # Run code linters
make test   # Run all tests
```

This helps maintain code quality and stability.

---

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
