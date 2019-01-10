# Build stage
FROM golang:1.11 as builder

RUN mkdir -p "$GOPATH/src/github.com/toskatok/lg"
WORKDIR $GOPATH/src/github.com/toskatok/pm
ENV GO111MODULE=on

COPY . .
RUN go build -o /bin/app

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /lg /app/
ENTRYPOINT ["./lg"]
