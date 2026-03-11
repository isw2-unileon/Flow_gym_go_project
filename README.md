# FlowGym

FlowGym is a web application designed to help gym users maintain the flow of their workouts during busy hours. When a machine is occupied, the system recommends alternative exercises that target the same muscle group using available equipment.

The goal of FlowGym is to reduce waiting times, improve workout efficiency, and maintain training intensity even when the gym is crowded.

---

# Features

- Exercise recommendation system based on muscle group
- Machine availability tracking
- Web interface for submitting workout requests
- Backend developed in Go
- Structured architecture using models, services, repositories and handlers
- Database integration for storing exercises and machines

---

# Requirements

To run the project locally you need:

- Go 1.21+
- PostgreSQL or SQLite database
- Git

---

# Environment Setup

Create a `.env` file in the root of the project:

```
DATABASE_URL=postgres://user:password@host:port/flowgym
PORT=8080
```

Replace the values with your actual database configuration.

---

# Running the Backend

Install dependencies and start the server:

```
go mod download
go run main.go
```

The backend server will run at:

```
http://localhost:8080
```

You can test the server health endpoint:

```
http://localhost:8080/health
```

---

# Database Setup

Initial database tables are defined in:

```
data/seed.sql
```

Run the SQL file to create the required tables:

- exercises
- machines
- users

---

# Project Structure

```
Flow_gym_go_project/
├── README.md
├── main.go
├── go.mod
│
├── data/
│   └── seed.sql
│
├── database/
│   └── db.go
│
├── handlers/
│   └── health_handler.go
│
├── models/
│   ├── exercise.go
│   ├── machine.go
│   ├── recommendation.go
│   └── user.go
│
├── repository/
│   ├── exercise_repository.go
│   └── machine_repository.go
│
├── services/
│   └── recommendation_service.go
│
├── templates/
│   └── index.html
│
└── static/
    └── css/
        └── style.css
```

---

# Architecture

The application follows a layered architecture:

**Handlers**
- Manage HTTP requests and responses

**Services**
- Implement the business logic (exercise recommendation)

**Repositories**
- Handle database queries

**Models**
- Define the core data structures of the application

**Database**
- Manages database connections

---

# How the Recommendation System Works

1. The user selects an exercise that cannot be performed because the machine is occupied.
2. The system identifies the target muscle group.
3. It checks which machines are currently available.
4. The backend recommends an alternative exercise that targets the same muscle group.

Example:

```
Requested exercise: Bench Press
Target muscle: Chest
Available equipment: Cable Machine

Recommendation: Cable Fly
```

---

# Future Improvements

- AI-based exercise recommendation
- Real-time machine availability tracking
- User authentication system
- Workout history and analytics
- Mobile-friendly interface

---

# Team

FlowGym development team:

- Samuel García Calvo
- Christian Castro Casado
- Roberto García Martín