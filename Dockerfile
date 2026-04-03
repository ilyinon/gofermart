FROM golang:1.26.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gophermart ./cmd/gophermart/main.go

CMD ["./gophermart"]
