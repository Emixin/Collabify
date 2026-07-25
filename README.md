# Collabify

**Collabify** is a lightweight team collaboration web application built with Go.  
It lets users sign up, log in, view users, teams, and tasks — inspired by the same core ideas as TeamUps (users with character types, teams with leaders, and tasks).

This is a simpler, Go-based implementation using the Gin web framework and GORM with SQLite.

## Features

- **User Management**
  - Sign up with username, email, and password
  - Login page
  - User character types: `Leader`, `Supporter`, `Doer`, `Thinker`, `Connector` (and `NoType`)
  - Basic user scoring

- **Teams**
  - Teams with a designated leader
  - Many-to-many relationship between teams and members

- **Tasks**
  - Tasks linked to teams
  - Deadline field

- **Pages**
  - Homepage
  - Login / Signup
  - Users list
  - Teams list
  - Tasks list

## Tech Stack

| Layer         | Technology                  |
|---------------|-----------------------------|
| Language      | Go 1.26+                    |
| Web Framework | Gin                         |
| ORM           | GORM                        |
| Database      | SQLite                      |
| Templates     | HTML (Gin HTML renderer)    |
| Static files  | CSS                         |

## Project Structure

```
Collabify/
├── cmd/
│   └── collabify/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── database.go          # SQLite + GORM setup & migrations
│   ├── handlers/
│   │   └── handlers.go          # HTTP handlers (pages + forms)
│   └── models/
│       └── models.go            # User, Team, Task models
├── web/
│   ├── statics/
│   │   └── home.css
│   └── templates/
│       ├── home.html
│       ├── login.html
│       ├── signup.html
│       ├── users_list.html
│       ├── teams_list.html
│       └── tasks_list.html
├── data/
│   └── db.sqlite                # SQLite database
├── go.mod
├── go.sum
└── .gitignore
```

## Getting Started

### Prerequisites

- Go 1.22+ (project uses Go 1.26.3 in `go.mod`)
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Emixin/Collabify.git
cd Collabify
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Run the application

```bash
go run ./cmd/collabify
```

The server starts at: **http://localhost:8080**

### Available Routes

| Method | Path          | Description          |
|--------|---------------|----------------------|
| GET    | `/`           | Homepage             |
| GET    | `/login`      | Login page           |
| POST   | `/login`      | Handle login form    |
| GET    | `/signup`     | Signup page          |
| POST   | `/signup`     | Create new user      |
| GET    | `/users_list` | List all users       |
| GET    | `/teams_list` | List all teams       |
| GET    | `/tasks_list` | List all tasks       |

## Models Overview

### User
- `ID`, `Username`, `Email`
- `Type` — one of: Leader, Supporter, Doer, Thinker, Connector, NoType
- `Score`

### Team
- `ID`, `Name`
- `Leader` (belongs to User)
- `Members` (many-to-many with User)

### Task
- `ID`, `Name`
- `Team` (belongs to Team)
- `Deadline` (string)

## Notes

- Authentication is currently basic (form handling exists, but session/auth middleware is not fully implemented yet in main branch).
- Passwords are accepted on signup/login but are **not yet hashed or stored** — this is improved in dev branch.
- Database is auto-migrated on startup via GORM.
- The SQLite file lives in `data/db.sqlite`.

## License

This project is currently unlicensed.  
Feel free to contact the author for collaboration.

---

Made by [Emixin](https://github.com/Emixin)
