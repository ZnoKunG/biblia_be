FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /bin/app ./cmd/api/*.go

FROM debian:buster-slim

COPY --from=build /bin/app /bin

COPY .env.* /bin

EXPOSE 3000

CMD ["/bin/app"]