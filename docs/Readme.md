# learn-share-backend

## For start project

> [!Note]
> If you want to start the project on your local PC, you should clone from the `local-start` branch

**Requirements:**

* [Git](https://git-scm.com/)
* [Docker](https://www.docker.com/)

**Steps:**

1. `git clone https://github.com/LearnShareApp/learn-share-backend.git`
    * or `git clone https://github.com/LearnShareApp/learn-share-backend.git -b local-start`
2. start docker on your PC
3. in terminal go to project root (learn-share-backend/)
    * enter command `docker-compose up -d --build`
4. for stop backend by command `docker-compose down`
5. after first run you can use `docker-compose up -d` command


## How to use

* after starting project you can open swagger and try api: http://adoe.ru:81
    * or http://localhost:81 for local start

## The main technologies which i use in this project

* [Go](https://go.dev/) - programming language
* [postgresql](https://www.postgresql.org/) - database
* [Docker](https://www.docker.com/) - containerization
* [chi](https://github.com/go-chi/chi) - http router
* [zap](https://github.com/uber-go/zap) - logging
* [swaggo](https://github.com/swaggo/swag) - swagger documentation
* [JWT](https://jwt.io/) - authentication
* [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - password hashing
* [sqlx](https://github.com/jmoiron/sqlx) - database communication
* [livekit](https://livekit.io/) - video communication (generating tokens for frontend)

## Contact

for questions contact with me in [telegram](https://t.me/Ruslan20007) or by email: ruslanrbb8@gmail.com



<!-- ## Learn-Share Backend Documentation


### Project Structure

The project follows a clean architecture pattern with clear separation of concerns:


1. **cmd/main** - Entry point with graceful shutdown

2. **internal/transport/rest** - REST API implementation using chi router

3. **internal/use_cases** - Business logic organized by domain areas

4. **pkg/logger** - Centralized logging with zap

5. **internal/config** - Configuration management


### Key Libraries

- **chi** - Lightweight HTTP router

- **zap** - High-performance logging

- **swaggo** - Swagger documentation generation

- **docker** - Containerization and deployment


### Use Case Structure

The use cases are organized by domain areas:


```go

internal/use_cases/

├── auth/            # Authentication

│   ├── login/       # User login

│   └── registration/ # User registration

├── categories/      # Category management

├── lessons/         # Lesson operations

│   ├── approve/     # Lesson approval

│   ├── book/        # Lesson booking

│   ├── cancel/      # Lesson cancellation

│   ├── finish/      # Lesson completion

│   ├── join/        # Lesson joining

│   ├── start/       # Lesson start

│   └── get/         # Lesson retrieval

├── schedules/       # Schedule management

│   ├── add_time/    # Add schedule time

│   └── get_times/   # Get available times

└── teachers/        # Teacher operations

    ├── add_skill/   # Add teacher skill

    ├── become/      # Become a teacher

    └── get/         # Teacher data retrieval

```


### Database Schema

Key tables and relationships:


1. **Users** - Core user information

2. **Teachers** - Teacher-specific data

3. **Lessons** - Lesson scheduling and management

4. **Categories** - Lesson categories

5. **Schedules** - Teacher availability

6. **Skills** - Teacher skills


Relationships:

- One-to-many between Teachers and Lessons

- Many-to-many between Teachers and Skills

- One-to-many between Categories and Lessons


### API Endpoints

The API is organized into these main routes:


```go

const (

    authRoute     = "/auth"      # Authentication

    userRoute     = "/user"      # User profile

    usersRoute    = "/users"     # Public user data

    teacherRoute  = "/teacher"   # Teacher operations

    teachersRoute = "/teachers"  # Public teacher data

    lessonRoute   = "/lesson"    # Lesson operations

    lessonsRoute  = "/lessons"   # Lesson management

    apiRoute      = "/api"       # Base API path

)

```


### Error Handling

The API uses consistent error handling:

- Standardized error responses

- Proper HTTP status codes

- Detailed error messages in development

- Secure error messages in production


### Deployment

The project uses Docker for containerization:

- `docker-compose up -d --build` for initial setup

- `docker-compose up -d` for subsequent runs

- `docker-compose down` to stop services


### Contact

For support, contact Ruslan via [Telegram](https://t.me/Ruslan20007) or email at ruslanrbb8@gmail.com.


This documentation provides a concise overview of the Learn-Share backend architecture and key components. For detailed API specifications, refer to the Swagger documentation at `http://localhost:81`. -->
