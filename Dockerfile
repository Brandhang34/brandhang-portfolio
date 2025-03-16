FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "./main" ]
