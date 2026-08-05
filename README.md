# Collabify

**Collabify** is a team collaboration web application built with Go.
Users can sign up, log in, create and manage teams, create and complete tasks, and view personalized dashboards — inspired by the same core ideas as TeamUps (users with character types, teams with leaders, and tasks).

Built with the Gin web framework, GORM, and SQLite.

## Features

- **Authentication**
  - Sign up with username, email, and password
  - Login / Logout
  - Session-based auth (cookie sessions via `gin-contrib/sessions`)
  - Passwords hashed with **bcrypt**
  - Account deletion with confirmation phrase

- **User Profiles**
  - Character types: `Leader`, `Supporter`, `Doer`, `Thinker`, `Connector`, `NoType`
  - Score + score count (running average support)
  - Dashboard with personal overview

- **Teams**
  - Create teams (with a designated leader)
  - Delete teams (leader only)
  - Many-to-many relationship between teams and members
  - Helper methods: `AddMember`, `RemoveMember`, `ChangeLeader`

- **Tasks**
  - Create tasks linked to a team
  - Delete tasks (leader only)
  - Status: `Pending` / `Completed`
  - Mark tasks as completed
  - Deadline field + `RenewDeadline` helper
  - Task list filtered to the current user’s teams

- **Other**
  - Homepage shows task counts (total + pending) for logged-in users
  - Users list, Teams list
  - Basic error handling utilities
  - Unit test helpers (in-memory SQLite)

## Tech Stack

| Layer           | Technology                          |
|-----------------|-------------------------------------|
| Language        | Go 1.26+                            |
| Web Framework   | Gin                                 |
| Sessions        | gin-contrib/sessions (cookie store) |
| Password Hashing| golang.org/x/crypto/bcrypt          |
| Config          | joho/godotenv                       |
| ORM             | GORM                                |
| Database        | SQLite                              |
| Templates       | HTML (Gin HTML renderer)            |
| Static files    | CSS                                 |

## Project Structure

```
Collabify/
├── cmd/
│   └── collabify/
│       └── main.go                 # Entry point, routes, session setup
├── internal/
│   ├── database/
│   │   └── database.go             # SQLite + GORM init & migrations
│   ├── handlers/
│   │   ├── handlers.go             # All HTTP handlers
│   │   └── handlers_test.go        # Tests
│   ├── middlewares/
│   │   └── middlewares.go
│   ├── models/
│   │   └── models.go               # User, Team, Task + methods
│   └── utils/
│       └── utils.go                # ErrorCatcher, UserTeamIDs, test helpers
├── web/
│   ├── statics/
│   │   ├── dashboard.css
│   │   ├── home.css
│   │   └── login.css
│   └── templates/
│       ├── home.html
│       ├── login.html
│       ├── signup.html
│       ├── dashboard.html
│       ├── users_list.html
│       ├── teams_list.html
│       ├── tasks_list.html
│       ├── create_team.html
│       ├── delete_team.html
│       ├── create_task.html
│       ├── delete_task.html
│       └── delete_account.html
├── data/
│   └── db.sqlite                   # SQLite database (created on first run)
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22+ (project uses Go 1.26.3)
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Emixin/Collabify.git
cd Collabify
```

### 2. Environment variables

Create a `.env` file in the project root:

```env
SESSION_SECRET=your-long-random-secret-here
```

> The app will panic on startup if `SESSION_SECRET` is missing.

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the application

```bash
go run ./cmd/collabify
```

Server starts at: **http://localhost:8080**

The SQLite database is created automatically at `data/db.sqlite` on first run (via GORM AutoMigrate).

## Routes

| Method(s)   | Path              | Description                          |
|-------------|-------------------|--------------------------------------|
| GET         | `/`               | Homepage (task summary if logged in) |
| GET / POST  | `/login`          | Login                                |
| GET / POST  | `/signup`         | Sign up                              |
| ANY         | `/logout`         | Logout                               |
| ANY         | `/dashboard`      | User dashboard                       |
| ANY         | `/delete_account` | Delete own account                   |
| GET         | `/users_list`     | List all users                       |
| GET         | `/teams_list`     | List all teams                       |
| GET / POST  | `/tasks_list`     | List user’s tasks / mark completed   |
| GET / POST  | `/create_team`    | Create a team                        |
| GET / POST  | `/delete_team`    | Delete a team (leader only)          |
| GET / POST  | `/create_task`    | Create a task                        |
| GET / POST  | `/delete_task`    | Delete a task (leader only)          |

## Models Overview

### User
- `ID`, `Username`, `PasswordHash`, `Email`
- `Type` — Leader / Supporter / Doer / Thinker / Connector / NoType
- `Score`, `ScoreCount`
- Methods: `UpdateUserAverageScore`, `UpdateUserType`

### Team
- `ID`, `Name`
- `Leader` (belongs to User)
- `Members` (many-to-many with User)
- Methods: `AddMember`, `RemoveMember`, `ChangeLeader`

### Task
- `ID`, `Name`, `Deadline`
- `Team` (belongs to Team)
- `Status` — Pending / Completed
- Methods: `RenewDeadline`, `MarkAsCompleted`

## License

This project is currently unlicensed.  
Feel free to contact the author for collaboration.

---

Made by [Emixin](https://github.com/Emixin)
