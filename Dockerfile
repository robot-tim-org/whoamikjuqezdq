FROM golang:alpine3.6 AS binary
ADD . /app
WORKDIR /app
RUN apk update && apk add git && go get gopkg.in/sirupsen/logrus.v1 && go get gopkg.in/spf13/viper.v1

RUN go vet .
RUN GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -a --installsuffix cgo -ldflags="-s -w" -o whoami .

FROM alpine:3.6
WORKDIR /app
ENV TEST_ENV="this is a test env variable"
COPY --from=binary /app/whoami /app
CMD ["/app/whoami", "-port", "80"]
EXPOSE 80
