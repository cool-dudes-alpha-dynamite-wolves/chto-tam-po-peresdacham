FROM golang:alpine
WORKDIR /chto-tam-po-peresdacham
COPY . .
RUN go mod download
EXPOSE 8081
ENTRYPOINT ["go", "run", "./cmd/main.go"]
