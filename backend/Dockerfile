FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main ./cmd/api/main.go

EXPOSE 8080

# Run the compiled binary
CMD ["./main"]
