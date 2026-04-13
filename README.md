# AICare Klondike API

## System Requirement

- Go 1.21.4
- Docker Compose

Following requirement provided in docker compose : 

- PostgreSQL
- Redis

## How to Install 

1. Clone project 
2. Copy ``.env.example`` to ``.env``
3. Adjust configuration in ``.env``
4. Run
   ```
   $ make install
   ```

## Build App

Once installation success run : 

```bash
// Dashboard API
$ make build  

// Private API
$ make build-private  
```

or build manually 

```bash
// Dashboard API
$ swag init --dir ./src/infrastructure/http/routes/dashboard --parseDependency true 
$ go build -o ./bin/app

// Private API
$ swag init --dir ./src/infrastructure/http/routes/private --parseDependency true 
$ go build -o ./bin/app
```

## Database Migration

Requirement : [sql-migrate](https://github.com/rubenv/sql-migrate)

Migration up 

```bash
$ sql-migrate up
```

Migration down 

```bash
$ sql-migrate up
```

### Generate Migration File

```bash
$ sql-migrate new <migration_name>
```

## Unit Test

```bash
$ go test ./... --cover -coverprofile=coverage.out
```
### Generate Mock
Requirement : [Mockery](https://github.com/vektra/mockery)

Run command

```bash
$ mockery --dir <interface dir> --output <output dir> --name <interface name>
```

Example

```bash
$ mockery --dir ./src/domain/ --output ./src/domain/mocks --name AccountRepository
```


## Run

By Terminal

```bash
// Dashboard API
$ make run

// Private API
$ make run-private
```

By Air In Terminal (Development only)

```bash
// Dashboard API
$ air

// Private API
$ air -c .air.private.toml
```