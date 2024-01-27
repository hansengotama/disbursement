# Disbursement Service

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

## Study Case
- Topic: The user has a balance in the application wallet and the balance wants to be
  disbursed.
  - Write code in Golang 
  - Write API (only 1 endpoint) for disbursement case only 
  - User data and balances can be stored as hard coded or database

### Implementation Criteria
To address the study case, certain criteria have been established:
- **Currency Restriction:** Disbursement is limited to a single currency (IDR).
- **Concurrency Control:** Users are prevented from triggering multiple disbursement requests simultaneously. This measure is implemented to avoid potential race conditions.
- **Admin Fee:** Disbursement incurs an admin fee, and users are required to cover this additional cost.
- **Payment Platform Options:** Users can register multiple payment platforms and choose the preferred one for their disbursement needs.

## Test the endpoint
```
  curl --location 'localhost:3000/disbursements/request' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9' \
--data '{
    "amount": 9800,
    "disbursement_account_guid": "6c2db1f4-0e63-4f47-8f77-2ac99acdbbc7"
}'
```

# Unit Test (Service Layer)
```bash
go test ./internal/service/disbursement --cover
```
