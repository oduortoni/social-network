# backend/Dockerfile
FROM golang:1.23.4-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 9000

CMD ["./main"]