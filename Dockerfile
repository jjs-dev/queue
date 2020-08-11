FROM golang:latest

WORKDIR /app

COPY COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]