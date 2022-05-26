FROM golang:1.17

WORKDIR /app

COPY . .

RUN go build -o forum ./cmd/web

ENV PORT 3000

EXPOSE $PORT

VOLUME [ "/app/data" ]

CMD ["./forum"]