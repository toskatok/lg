# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/mqttlg
RUN apk update && apk add git
RUN cd $GOPATH/src/github.com/aiotrc/mqttlg/ && go get -v && go build -v -o /mqttlg

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /mqttlg /app/
ENTRYPOINT ["./mqttlg"]
