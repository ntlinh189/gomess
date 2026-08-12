# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gomess-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /gomess-cronjob ./cmd/cronjob

# ---- Runtime stage ----
FROM alpine:3.20 AS api
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /gomess-api .
RUN adduser -D -u 1000 appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["./gomess-api"]

# ---- Runtime: Worker ----
FROM alpine:3.20 AS cronjob
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /gomess-cronjob .
RUN adduser -D -u 1000 appuser
USER appuser
ENTRYPOINT ["./gomess-cronjob"]