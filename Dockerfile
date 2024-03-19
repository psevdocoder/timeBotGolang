FROM golang:1.21-alpine AS builder
LABEL authors="Maksim"

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
ADD go.mod .
COPY . .

RUN go build -o BotApp ./cmd/main.go

FROM alpine:3.19

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata

ENV TZ=Europe/Moscow
VOLUME /app/config
VOLUME /app/logs

WORKDIR /app
COPY --from=builder /app/BotApp ./BotApp


CMD ["/app/BotApp"]