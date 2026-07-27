FROM golang:1.26 AS builder
WORKDIR /app

COPY pkg ./pkg
COPY services/order/go.mod services/order/go.sum ./services/order/
COPY services/user/go.mod services/user/go.sum ./services/user/
COPY services/product/go.mod services/product/go.sum ./services/product/

RUN go -C services/order mod download && \
    go -C services/user mod download && \
    go -C services/product mod download

COPY services ./services

RUN CGO_ENABLED=0 GOOS=linux go -C services/order build -o /app/order ./cmd && \
    CGO_ENABLED=0 GOOS=linux go -C services/user build -o /app/user ./cmd && \
    CGO_ENABLED=0 GOOS=linux go -C services/product build -o /app/product ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/order /app/user /app/product ./
