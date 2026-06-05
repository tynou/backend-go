FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

ENV GOWORK=off

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/payment/go.mod services/payment/go.sum ./services/payment/

RUN cd services/payment && go mod download

COPY pkg/ ./pkg/
COPY services/payment/ ./services/payment/

WORKDIR /build/services/payment

RUN CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./cmd/


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /payment-service ./payment-service
COPY services/payment/db/migrations ./db/migrations

EXPOSE 8081

CMD ["./payment-service"]