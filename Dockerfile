FROM golang:1.24-alpine AS build

ARG APP_DIR=server

ENV CGO_ENABLED=0

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go generate ./...
RUN go build -o main ./cmd/api/*.go

##############################################################################################

FROM debian:buster-slim

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/cmd/generated/docs ./generated/docs

COPY .env.* .

EXPOSE 3000

CMD ["./main"]