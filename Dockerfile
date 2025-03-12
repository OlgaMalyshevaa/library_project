FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o apiserver cmd/apiserver/main.go

CMD ["./apiserver"]

EXPOSE 8080
