FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .env
RUN go build -o main ./src/cmd/public
CMD ["./main"]