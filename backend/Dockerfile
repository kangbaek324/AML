FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o aml main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/aml .

EXPOSE 8080
CMD ["./aml"]
