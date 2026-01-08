# =========================
# BUILD STAGE
# =========================
FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on

RUN apk add --no-cache git

# ✅ FIX: chỉ copy go.mod
COPY go.mod ./
RUN go mod download

COPY . .

# build đúng entrypoint
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# =========================
# RUN STAGE
# =========================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
