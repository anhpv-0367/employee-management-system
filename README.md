### Create file go.mod
```
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine \
  go mod init app
```

### Create file go.sum
```
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine \
  go get github.com/lib/pq
```

### migration database

```
export $(grep -v '^#' .env | xargs)
```

```
docker run --rm \
  --network container:postgres_db \
  -v "$PWD/migrations":/migrations \
  migrate/migrate \
  -path /migrations \
  -database "$DATABASE_URL" \
  up
```

#### check db:

```
docker compose exec db psql -U postgres -d employee_db

\dt
```

### CURL Example

- Create department:

```
curl -i -X POST http://localhost:8080/departments \   
  -H "Content-Type: application/json" \
  -d '{
    "name": "IT"
  }'
```

- Create employee:

```
curl -i -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nguyen Van A",
    "email": "a@example.com",
    "departmentId": 1
  }'
```

- Get Employee Detail

```
curl -i http://localhost:8080/employees/1
```

- GET /employees

```
curl --location 'http://localhost:8080/employees?limit=1&offset=2&departmentId=1'
```

- PUT /employees/:id

```
curl --location --request PUT 'http://localhost:8080/employees/11' \
--header 'Content-Type: application/json' \
--data-raw '{
    "position": "Leaderx",
    "age": 12,
    "salary": 12.3,
    "name": "Nguyen van fix",
    "email": "D@example.com",
    "DepartmentID": 1
  }'
```

- GET /employees/:id

```
curl --location --request DELETE 'http://localhost:8080/employees/12' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Nguyen Van B",
    "email": "b@example.com",
    "department_id": 1
  }'
```