# library

REST APIs of a book library.

## Tech
- Gorilla `mux` + `handlers` (CORS)
- `zerolog` logging
- Config with `viper` (YAML + env overrides)
- Postgres via `pgx`, migrations via `goose`
- Validation via `validator.v10`
- Tests with `testify/require` (+ integration via `testcontainers-go`)

## Quick start
- Docker Compose:
```bash
make up
# http://localhost:8080/healthz
# http://localhost:8080/v1/books
```
## Config
- Default file: `config/config.yaml`
- Env overrides: prefix `LIB_` (e.g. `LIB_DB_HOST=db`)

## Make targets
- `build`, `run`, `test`, `up`, `down`, `logs`, `docker-build`, `migration-create`



