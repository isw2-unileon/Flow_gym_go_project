# FlowGym: Technical Documentation

## 1. Project Organization and Requirements

The core functional requirements implemented are:
- User registration and login (authentication).
- Management of user and administrator roles.
- Real-time querying of machine availability.
- Occupying and releasing machines.
- Automatic machine release after 15 minutes of use.
- Recommendation of alternative exercises based on muscle groups.
- CRUD operations for personalized workout routines.

Project organization was managed using **GitHub Projects with a Kanban board**. The features were planned through a structured backlog, and each task was assigned to a specific team member to ensure traceability of individual contributions.

## 2. Architecture and Design

The application follows a layered architecture, clearly separating the UI, business logic, and data access.

| Technology | Usage |
| :--- | :--- |
| Go (Golang) | Backend development |
| PostgreSQL | Relational Database |
| HTML, CSS, JS | Frontend development |
| GitHub | Version control and CI/CD |
| Render | Production Deployment |

**System Flow:**
Client (Browser) -> Frontend (HTML/CSS/JS) -> Backend Server (Go) -> PostgreSQL Database.

## 3. Development and Quality Assurance

Development followed industry best practices, including Git version control, integration via **Pull Requests**, and separation of concerns.

A major focus was placed on software quality. We implemented a **Test-Driven Development (TDD)** approach, achieving **100% Test Coverage** in the critical backend layers (handlers, repositories, and services). We utilized the `sqlmock` library to simulate complex database transactions, ensuring robust error handling without compromising real data.

## 4. CI/CD and Deployment Workflow

The project followed the **Trunk Based Development** methodology. Features were developed in short-lived branches, peer-reviewed via Pull Requests, and frequently integrated into the `main` branch. 

We implemented **CI/CD pipelines using GitHub Actions** to automate the execution of our Go test suite, ensuring code integrity before every merge. 

The application is deployed on Render, hosting both the Go server and the PostgreSQL database.