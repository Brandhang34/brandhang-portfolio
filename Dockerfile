# Fetch
FROM golang:1.23-alpine3.21 AS fetch-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
WORKDIR /app
COPY --chown=65532:65532 . .
RUN ["templ", "generate"]

# Build
FROM golang:1.23-alpine3.21 AS build-stage
WORKDIR /app
COPY --from=generate-stage /app .
RUN go build -o main cmd/api/main.go

# Test
# Not in the default build graph (nothing downstream depends on it) -- run it
# explicitly with `docker build --target test-stage .`. CI is the real gate.
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy
FROM alpine:latest
WORKDIR /app
COPY --from=build-stage /app/main .
EXPOSE 8080
ENV PORT=8080
CMD [ "./main" ]
