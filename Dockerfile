FROM golang:1.21.4 as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./photo-app .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/photo-app /app/.env ./

EXPOSE ${APP_PORT}

CMD "./photo-app"