FROM golang:1.14-alpine
COPY .  /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/chatserver ./cmd

EXPOSE 9000
CMD ["/bin/chatserver"]