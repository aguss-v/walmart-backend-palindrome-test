FROM  golang:1.17.1-stretch

WORKDIR /app

COPY . .

RUN `go mod tidy \
  && go build src/cmd/main.go`

CMD ["./main"]
