FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/analytics-worker ./cmd/analytics-worker


FROM alpine:3.22 AS api

WORKDIR /app

COPY --from=builder /bin/api ./api
COPY --from=builder /app/seed .

EXPOSE 8080

CMD ["./api"]


FROM alpine:3.22 AS analytics-worker

WORKDIR /app

COPY --from=builder /bin/analytics-worker ./analytics-worker

CMD ["./analytics-worker"]