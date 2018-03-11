FROM golang:1.9
RUN apt-get update

RUN go get github.com/gorilla/mux && \
    go get github.com/gorilla/handlers && \
    go get -u github.com/go-sql-driver/mysql && \
    go get golang.org/x/crypto/bcrypt && \
    go get github.com/nu7hatch/gouuid

COPY . /go/src/eventplanner
WORKDIR /go/src/eventplanner/src/server
RUN ["./gotest.sh"]
COPY . /go/src/eventplanner
WORKDIR /go/src/eventplanner
RUN ["./build.sh"]
COPY ./entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
