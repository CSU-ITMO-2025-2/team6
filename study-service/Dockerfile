FROM golang:1.24.6-alpine AS deps

WORKDIR /app

COPY study-service/go.mod study-service/go.sum ./study-service/
COPY local-lib/database/go.mod local-lib/database/go.sum ./local-lib/database/
COPY local-lib/queue/go.mod local-lib/queue/go.sum ./local-lib/queue/

WORKDIR /app/study-service
RUN go mod download

FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY --from=deps /go/pkg/mod /go/pkg/mod

COPY study-service/ ./study-service/
COPY local-lib/ ./local-lib/

WORKDIR /app/study-service

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o /app/main cmd/grpc_server/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 5000

CMD ["./main"]