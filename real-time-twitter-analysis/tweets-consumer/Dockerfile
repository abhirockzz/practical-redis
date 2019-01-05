FROM golang as build-stage
WORKDIR /go/
RUN go get github.com/go-redis/redis
COPY src/ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tweets-consumer

FROM scratch
COPY --from=build-stage /go/src/tweets-consumer /
CMD ["/tweets-consumer"]