FROM golang:1.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o forum ./cmd/web

FROM alpine:latest AS production
COPY --from=builder /app .
EXPOSE 5050
CMD ["./forum"]