FROM golang
WORKDIR /go/src/app
RUN go get github.com/go-redis/redis
COPY src/tweets-consumer.go /go/src/app
CMD ["go","run","tweets-consumer.go"]