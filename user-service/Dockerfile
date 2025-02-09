FROM golang:1.23 AS builder
WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8082

CMD ["/app/main"]