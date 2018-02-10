FROM golang:1.9
RUN apt-get update

RUN go get github.com/gorilla/mux && \
    go get github.com/gorilla/handlers && \
    apt-get install -y vim

COPY . /go/src/event-planner
WORKDIR /go/src/event-planner
RUN ["./build.sh"]
COPY ./entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
EXPOSE 8080
