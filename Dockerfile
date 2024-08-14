FROM build_image AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -a -o server ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .
COPY ./config.yml .
CMD ["./server"]
