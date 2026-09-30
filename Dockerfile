FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api \
    && CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate \
    && CGO_ENABLED=0 GOOS=linux go build -o /out/seed ./cmd/seed \
    && CGO_ENABLED=0 GOOS=linux go build -o /out/billing-cron ./cmd/billing-cron

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/api ./api
COPY --from=builder /out/migrate ./migrate
COPY --from=builder /out/seed ./seed
COPY --from=builder /out/billing-cron ./billing-cron

EXPOSE 8080

ENTRYPOINT ["./api"]
