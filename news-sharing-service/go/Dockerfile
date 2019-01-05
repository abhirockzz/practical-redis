FROM golang as build-stage
WORKDIR /go/
RUN go get -u github.com/go-redis/redis && go get -u github.com/gin-gonic/gin
COPY src/ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o news-sharing-app

FROM scratch
COPY --from=build-stage /go/src/news-sharing-app /
CMD ["/news-sharing-app"]
