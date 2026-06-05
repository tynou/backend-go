FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

ENV GOWORK=off

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/

RUN cd services/gateway && go mod download

COPY pkg/ ./pkg/
COPY services/gateway/ ./services/gateway/

WORKDIR /build/services/gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-service ./cmd/


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /gateway-service ./gateway-service

EXPOSE 8084

CMD ["./gateway-service"]