FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

ENV GOWORK=off

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/

RUN cd services/auth && go mod download

COPY pkg/ ./pkg/
COPY services/auth/ ./services/auth/

WORKDIR /build/services/auth

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /auth-service ./auth-service
COPY services/auth/db/migrations ./db/migrations

EXPOSE 8081

CMD ["./auth-service"]