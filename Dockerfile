FROM golang:1.9
RUN apt-get update

RUN go get github.com/gorilla/mux && \
    go get github.com/gorilla/handlers && \
    go get -u github.com/go-sql-driver/mysql && \
    go get golang.org/x/crypto/bcrypt && \
    go get github.com/nu7hatch/gouuid

## UNCOMMENT TO RUN TESTS.
## Don't use docker-compose up to run tests.
## Use ./runtest.sh which will start the needed
## mysql db first.
COPY . /go/src/eventplanner
WORKDIR /go/src/eventplanner/src/server
RUN ["./gotest.sh"]

## UNCOMMENT TO RUN APP.
#COPY . /go/src/eventplanner
#WORKDIR /go/src/eventplanner
#RUN ["./build.sh"]
#COPY ./entrypoint.sh /
#ENTRYPOINT ["/entrypoint.sh"]
