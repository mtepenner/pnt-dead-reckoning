FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go mod download
RUN go build -o /navigator ./cmd/navigator

FROM alpine:3.20
COPY --from=build /navigator /usr/local/bin/navigator

EXPOSE 8081 9101
ENTRYPOINT ["/usr/local/bin/navigator"]
