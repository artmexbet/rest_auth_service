FROM build_image
WORKDIR /app
COPY . .
RUN GOOS=linux go build -a -o server ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app
COPY --form=builder /app/server .
COPY ./config.yml .
CMD ["./server"]
