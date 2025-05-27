# bitebattle Backend

bitebattle Backend is a Go-based backend application designed to handle various functionalities related to the bitebattle project. This document provides an overview of the project, setup instructions, and usage guidelines.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Installation

To get started with the bitebattle Backend, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/bitebattle-backend.git
   ```

2. Navigate to the project directory:
   ```
   cd bitebattle-backend
   ```

3. Install the dependencies:
   ```
   go mod tidy
   ```

## Usage

To run the application, use the following command:
```
go run cmd/main.go
```

The server will start and listen on the specified port. You can then access the API endpoints as defined in the application.

## Project Structure

```
bitebattle-backend
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── handler
│   │   └── handler.go   # HTTP request handling logic
│   ├── service
│   │   └── service.go   # Business logic and services
│   └── model
│       └── model.go     # Data structures and models
├── go.mod                # Module dependencies
├── go.sum                # Dependency checksums
└── README.md             # Project documentation
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

## UML 

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
                    | status (open/closed) |
                    | created_at           |
                    +----------------------+

                             ▲
                             │
                    +----------------------+
                    |     poll_votes       |
                    +----------------------+
                    | id (PK)              |
                    | poll_id (FK→polls)   |
                    | user_id (FK→users)   |
                    | restaurant_id (FK→restaurants) |
                    | voted_at             |
                    +----------------------+
