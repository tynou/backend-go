FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.work ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY services/billing/go.mod services/billing/go.sum ./services/billing/
COPY services/payment/go.mod services/payment/go.sum ./services/payment/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/

RUN go work sync

COPY pkg/ ./pkg/
COPY services/auth/ ./services/auth/

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./services/auth/cmd/


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /auth-service ./auth-service

COPY services/auth/db/migrations ./db/migrations

EXPOSE 8081

CMD ["./auth-service"]