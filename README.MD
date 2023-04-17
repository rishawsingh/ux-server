Run docker compose to setup dependencies
```bash
docker compose up -d
```

Setup environment variable
```bash
DB_HOST=localhost
DB_NAME=sample
DB_PORT=5433
DB_USER_PASSWORD=local
DB_USERNAME=local
SERVICE_NAME=sample
```

Linting
```bash
# install golangci linter then
make lint
```