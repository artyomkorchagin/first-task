FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/order-service ./cmd/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/order-service /app/order-service
COPY  ./migrations /app/migrations
#COPY ./docs /app/docs
COPY ./.env /app/.env
EXPOSE 3000
CMD ["/app/order-service"]