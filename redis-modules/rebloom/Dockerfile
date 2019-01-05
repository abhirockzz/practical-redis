FROM golang:alpine as build-stage
WORKDIR /go/
RUN apk --no-cache add ca-certificates git
RUN go get -u github.com/go-redis/redis && go get -u github.com/gorilla/mux
COPY src/ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o articles-recommendation-service

FROM scratch
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /go/src/ /
CMD ["/articles-recommendation-service"]