FROM golang as build-stage
WORKDIR /go/
RUN go get github.com/go-redis/redis && go get github.com/gin-gonic/gin
COPY src/ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tweets-stats

FROM scratch
COPY --from=build-stage /go/src/tweets-stats /
CMD ["/tweets-stats"]