# BiteBattle

BiteBattle is a collaborative decision-making app for groups to vote on restaurants and settle food debates. It consists of an iOS client and a Go-based backend server.

---

## iOS App

The BiteBattle iOS app provides a user-friendly interface for creating groups, joining polls, voting on restaurant options, and managing your account. It is built with SwiftUI and communicates with the backend via RESTful APIs.

---

## Server (Backend)

The backend is written in Go and handles user authentication, group management, poll creation, voting, and notifications. It uses PostgreSQL for data storage and exposes RESTful endpoints for the mobile app.

For more details, see the [bitebattle-backend README](bitebattle-backend/README.md).
