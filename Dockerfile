FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o reddit-app .

FROM alpine

WORKDIR /app
COPY --from=builder /app/reddit-app .

EXPOSE 8080

CMD ["./reddit-app"]
