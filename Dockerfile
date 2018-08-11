# Build stage
FROM golang:alpine AS build-env
COPY . $GOPATH/src/github.com/I1820/lg
RUN apk --no-cache add git
WORKDIR $GOPATH/src/github.com/I1820/lg
RUN go get -v && go build -v -o /lg

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /lg /app/
ENTRYPOINT ["./lg"]
