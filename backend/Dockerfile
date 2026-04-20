FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o app ./cmd/app

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["/app/app"]
