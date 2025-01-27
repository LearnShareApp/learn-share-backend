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
2. create `.env` file in project root and fill it with your data (see `.env.example`)
3. start docker on your PC
4. in terminal go to project root (learn-share-backend/)
    * enter command `docker-compose up -d --build`
5. for stop backend by command `docker-compose down`
6. after first run you can use `docker-compose up -d` command


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

## Project structure

### Code
#### Main structure
```go
cmd/
├── main/main.go               // start project, load config, graceful shutdown
internal/
├── application/application.go // application, app managment: conect to db, logger, create/run/stop server
├── config/config.go           // load, validate config from env
├── entities/...               // entities for db and business logic
├── errors/...                 // errors in business logic, repo and transport
├── jsonutils/...              // json utils for transport
├── repository/...             // db repository: init, queries, transactions
├── service
│   ├── jwt/service.go            // jwt: generate, validate
│   └── livekit/service.go        // livekit: generate meeting tokens
├── transport/rest/
│   ├── server.go              // rest api, routing, adding handlers & middlewares
│   └── middlewares/...        // middlewares: auth (jwt), cors, logger
├── use_cases/...              // business logic cases: endpoint, service, repo, type
pkg/
├── db/postgres/postgres.go    // postgres: connect
├── hasher/hasher.go           // hash: bcrypt
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
* `nginx` - reverse proxy
* `app-1` - backend application
## Contact

for questions contact with me in [telegram](https://t.me/Ruslan20007) or by email: ruslanrbb8@gmail.com