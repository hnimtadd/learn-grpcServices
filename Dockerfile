FROM golang:1.21.3-alpine3.18 AS BUILDER

WORKDIR /app
COPY . .
RUN go build -o main ./cmd/server/main.go

FROM alpine:3.18 AS BASE
WORKDIR /app
COPY --from=BUILDER /app/main .
EXPOSE 8080

CMD ["/app/main"]
