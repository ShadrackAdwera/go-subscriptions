# Build Stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migrations

EXPOSE 5000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]