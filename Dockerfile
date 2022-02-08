
FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./


# build go app
RUN go mod download
RUN go build -o second ./cmd/main/app.go

CMD ["./cmd/main/second"]