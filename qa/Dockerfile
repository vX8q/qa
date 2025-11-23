FROM golang:1.24-alpine

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/qa ./cmd

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV PATH=$PATH:/root/go/bin

COPY ./migrations /app/migrations

CMD ["sh", "-c", "goose -dir /app/migrations postgres \"host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable\" up && /app/qa"]
