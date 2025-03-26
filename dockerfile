FROM golang:1.24.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY api/docs/v1/swagger.json /app/./api/docs/v1/swagger.json
EXPOSE 8080
CMD ["./main"]
