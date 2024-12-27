FROM golang:alpine
WORKDIR /chto-tam-po-peresdacham
COPY . .
RUN go mod download
EXPOSE ${{ secrets.PORT }}
ENTRYPOINT ["go", "run", "./cmd/main.go"]
