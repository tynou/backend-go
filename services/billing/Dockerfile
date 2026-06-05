FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

ENV GOWORK=off

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/billing/go.mod services/billing/go.sum ./services/billing/

RUN cd services/billing && go mod download

COPY pkg/ ./pkg/
COPY services/billing/ ./services/billing/

WORKDIR /build/services/billing

RUN CGO_ENABLED=0 GOOS=linux go build -o /billing-service ./cmd/


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /billing-service ./billing-service
COPY services/billing/db/migrations ./db/migrations

EXPOSE 8083

CMD ["./billing-service"]