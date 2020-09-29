FROM golang:1.15.2-alpine3.12 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o shortener.bin .

FROM alpine:3.12
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/shortener.bin /app/shortener
WORKDIR /app
ENTRYPOINT ["./shortener"]
CMD ["run", "--bind-address=:8080"]
