ARG GO_VERSION=1.17.5
FROM golang:$GO_VERSION as build

WORKDIR /bot
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app

FROM alpine:latest

COPY --from=build /bot/config /config
COPY --from=build /bot/app /app

ENTRYPOINT ["/app", "-config=/config/config.prod.json"]