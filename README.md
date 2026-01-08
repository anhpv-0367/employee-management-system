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