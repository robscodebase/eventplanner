FROM golang:1.9
RUN apt-get update

RUN go get github.com/gorilla/mux && \
    go get github.com/gorilla/handlers && \
    go get -u github.com/go-sql-driver/mysql

COPY . /go/src/eventplanner
WORKDIR /go/src/eventplanner
RUN ["./build.sh"]
COPY ./entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
