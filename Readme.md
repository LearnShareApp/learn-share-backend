# learn-share-backend

## To start the project

> [!Note]
> If you want to start the project on your local PC, you should clone from the `local-start` branch

**Requirements:**

* [Git](https://git-scm.com/)
* [Docker](https://www.docker.com/)

**Steps:**

1. `git clone https://github.com/LearnShareApp/learn-share-backend.git`
    * or `git clone https://github.com/LearnShareApp/learn-share-backend.git -b local-start`
2. create `.env` file in project root and fill it with your data (see `.env.example`)
3. start docker on your PC
4. in terminal go to project root (learn-share-backend/)
    * enter command `docker-compose up -d --build`
5. for stop backend by command `docker-compose down`
6. after first run you can use `docker-compose up -d` command


## How to use

After starting the project, you can open Swagger and try the API: http://adoe.ru:81/swagger/
(or http://localhost:81/swagger/ for local start) <br>

## Rest API handlers

You can see all handlers in my swagger: http://adoe.ru:81/swagger/index.html
(or http://localhost:81/swagger/index.html for local start)
But for sure I describe handlers also here:
Here's a concise tree-structured description of the API handlers:

``` go
API Handlers
├── Auth
│   ├── POST /auth/login             // Login with email and password
│   └── POST /auth/signup            // Register a new user (student)
│ 
├── Categories 
│   └── GET /categories              // Get list of all categories
│ 
├── Image 
│   └── GET /image                   // Get image by filename
│
├── Lessons
│   ├── POST /lesson                 // Add unconfirmed lesson (lesson request)
│   ├── GET /lessons                 // Get lessons for students
│   ├── GET /lessons/{id}            // Get lesson data by lesson's id
│   ├── GET /lessons/{id}/short-data // Get lesson really short data by lesson's id
│   ├── PUT /lessons/{id}/approve    // Approve lesson
│   ├── PUT /lessons/{id}/cancel     // Cancel lesson
│   ├── PUT /lessons/{id}/finish     // Finish lesson
│   ├── GET /lessons/{id}/join       // Join the lesson (generate meet token)
│   └── PUT /lessons/{id}/start      // Start lesson (generate meet token)
│
├── Reviews
│   └── POST /review                 // Create review
│
├── Teachers
│   ├── GET /teacher                 // Get teacher data (by JWT)
│   ├── POST /teacher                // Register as a teacher
│   ├── GET /teacher/lessons         // Get lessons for teachers
│   ├── GET /teacher/schedule        // Get times from schedule
│   ├── POST /teacher/schedule       // Add time to schedule
│   ├── POST /teacher/skill          // Register new skill for teacher
│   ├── GET /teachers/{id}/reviews   // Get all about teacher by his teacher's ID
│   ├── GET /teachers                // Get full teachers data
│   ├── GET /teachers/{id}           // Get teacher data by ID
│   └── GET /teachers/{id}/schedule  // Get times from schedule by ID
│ 
└── Users 
    ├── GET /user/profile            // Get user profile (by JWT)
    ├── PATCH /user/profile          // Edit user profile
    └── GET /users/{id}/profile      // Get user profile by UserID
```



## The main technologies which I use in this project

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
* [minio](https://min.io/) - object storage for images

## Project structure

### Code
#### Main structure
```
cmd/
├── main/main.go               // start project, load config, graceful shutdown
internal/
├── application/application.go // application, app managment: conect to db, logger, create/run/stop server
├── config/config.go           // load, validate config from env
├── entities/...               // entities for db and business logic
├── errors/...                 // errors in business logic, repo and transport
├── imgutils/...               // image utils for operations under image 
├── httputils/...              // http utils for transport layer
├── repository/...             // db repository: init, queries, transactions
├── service/
│   ├── jwt/service.go         // jwt: generate, validate
│   └── livekit/service.go     // livekit: generate meeting tokens
├── transport/rest/
│   ├── server.go              // rest api, routing, adding handlers & middlewares
│   └── middlewares/...        // middlewares: auth (jwt), cors, logger
├── use_cases/...              // business logic cases: endpoint, service, repo, type
pkg/
├── db/postgres/postgres.go    // postgres: connect
├── hasher/hasher.go           // hash: bcrypt
├── object_storage/...         // object storage: for saving images
└── logger/logger.go           // logger: zap, configurations
```
#### Use_case structure
Use cases grouped in by their 'domains' or 'features' (for example: auth, lessons, etc.).<br>
Every use case has 4 files (some of them have 3 files)
* `endpoint.go` - transport layer, wraps business logic and create http handler
* `service.go` - business logic, main logic of this case
* `repo.go` - db repository, functions for db communication in service
* `type.go` - mapping types: request, response of use case (marshal, unmarshal json). Some use_cases haven't this file.

### Docker
3 containers:
* `postgres` - database
* `minio` - object storage
* `nginx` - reverse proxy
* `app-1` - backend application

## Storing data

* The project utilizes PostgreSQL as the database management system.
* The repository handles table initialization and establishes relationships between entities, including tables for users, teachers, categories, skills, schedules, statuses, lessons, and reviews.
* Triggers are implemented to automatically update teacher and skill ratings when new reviews are added.
* The database is pre-populated with essential seeding data for categories and statuses to ensure proper functionality from the start.

* All methods for data operations are split into files by logical parts and use entity structs:
    * category
    * lesson
    * review
    * schedule
    * skill
    * status
    * teacher
    * user

## Contact

for questions contact with me in [telegram](https://t.me/Ruslan20007) or by email: ruslanrbb8@gmail.com
