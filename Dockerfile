FROM golang:1.9
RUN apt-get update
ENV NODE_VERSION 8.9.3
ENV NPM_VERSION 5.0.0
RUN curl -SLO "http://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz" \
    && tar -xzf "node-v$NODE_VERSION-linux-x64.tar.gz" -C /usr/local --strip-components=1 \
    && npm install -g npm@"$NPM_VERSION"
ENV PATH $PATH:/nodejs/bin

RUN go get github.com/gorilla/mux && \
    go get golang.org/x/net/context && \
    apt-get install -y vim

COPY . /go/src/event-planner
WORKDIR /go/src/event-planner/src/server
RUN ["./build.sh"]
COPY ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

