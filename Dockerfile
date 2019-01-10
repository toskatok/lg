# Build stage
FROM golang:1.11 as builder

RUN mkdir -p "$GOPATH/src/github.com/toskatok/lg"
WORKDIR $GOPATH/src/github.com/toskatok/pm
ENV GO111MODULE=on

COPY . .
RUN go build -o /bin/app

# Final stage
FROM alpine:latest

WORKDIR /bin

COPY --from=builder /bin/app .

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

CMD ["/bin/app"]
