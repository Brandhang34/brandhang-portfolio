# TEMPL_VERSION must match the templ runtime in go.mod. It was ':latest', which
# generated code calling functions the pinned library did not export
# (undefined: templ.JoinURLErrs). Bump both together.
ARG GO_VERSION=1.25
ARG TEMPL_VERSION=v0.3.1020

# Fetch
FROM golang:${GO_VERSION}-alpine AS fetch-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:${TEMPL_VERSION} AS generate-stage
WORKDIR /app
COPY --chown=65532:65532 . .
RUN ["templ", "generate"]

# Build
FROM golang:${GO_VERSION}-alpine AS build-stage
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
