# disbursement

## Stack:
- Golang (version 1.21.3)
- PostgreSQL
- Redis

## Prerequisite
To run the project, ensure you have the following prerequisites installed:

### Installation
- Install Go Migrate
  - Mac OS:
    ```bash
    brew install golang-migrate
    ```
  - Others: [Check here](https://www.freecodecamp.org/news/database-migration-golang-migrate/)


- Install Mockery 
    ```bash
    go install github.com/vektra/mockery/v2@v2.38.0
    ```

- Install direnv
  - Mac OS:
    ```bash
    brew install direnv
    ```
  - Linux: 
    ```bash
    sudo apt-get install direnv
    ```

### ENV
- Copy the `.env.example` file to a new file named `.env`
  ```bash
  cp .env.example .env
  ```
- Update the .env file with relevant configuration parameters. See [direnv documentation](https://direnv.net) for details.
  ```bash
  direnv allow
  ```

## Migration
We utilize [golang migrate](https://github.com/golang-migrate/migrate) for managing database migrations
### Migrate SQL Schema
Apply Migration (Up)

To create a new SQL schema, execute the following command:
```bash
make migrate-up
```
Revert Migrations (Down)

If you need to revert applied migrations, use the following command:

```bash
make migrate-rollback
```

### Add a New SQL Schema
To create a new SQL schema, execute the following command:
```
migrate create -ext sql -dir sql/schema/ -seq {{FILENAME}}
```

For instance, to create a schema named create_table_users, run:
```
migrate create -ext sql -dir sql/schema/ -seq create_table_users
```

This allows you to generate new migration files for your SQL schema changes.

### Seed
The seed command is intended to be executed during the initial setup and should be run only once.
```bash
go run . seed-please
```

# How to run
```bash
go run .
```

# Unit Test

