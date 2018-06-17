FROM golang
WORKDIR /go/
RUN go get github.com/go-redis/redis
RUN go get github.com/gorilla/mux
COPY src/ /go/src
CMD ["go","run","src/main.go"]