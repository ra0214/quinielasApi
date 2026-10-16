FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o quinielas .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/quinielas .
EXPOSE 8080
CMD ["./quinielas"]
