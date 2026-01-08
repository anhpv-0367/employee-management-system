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