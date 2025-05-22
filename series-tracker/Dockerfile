
FROM golang:1.23.8

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod tidy

# Install swag CLI (for Swagger docs)
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/air-verse/air@latest

# Build your app from the root (where main.go is)
#RUN go build -o main .


EXPOSE 8080
#CMD ["air", "-c", ".air.toml"]
