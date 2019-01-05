FROM golang as build-stage
WORKDIR /go/
RUN go get -u github.com/go-redis/redis && go get -u github.com/gin-gonic/gin && go get -u github.com/dghubble/oauth1 && go get -u github.com/dghubble/go-twitter/twitter
COPY src/ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tweets-producer

FROM scratch
COPY --from=build-stage /go/src/tweets-producer /
CMD ["/tweets-producer"]
