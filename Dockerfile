FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o /connect4 ./cmd/api

FROM gcr.io/distroless/base-debian12
COPY --from=builder /connect4 /connect4
EXPOSE 8080
ENTRYPOINT ["/connect4"]
