FROM ubuntu:14.04

RUN  apt-get update && \
     apt-get install -y curl && \
     apt-get install -y cmake && \
     apt-get install -y git && \
     rm -rf /var/lib/apt/lists/*

#Get go1.6
RUN curl https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz | tar -xvzf - -C /usr/local

#gopath
RUN mkdir go
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/revel/cmd/revel
RUN go get bitbucket.org/liamstask/goose/cmd/goose
RUN go get github.com/aws/aws-sdk-go
RUN go get upper.io/db.v2/postgresql

RUN mkdir -p /go/src/github.com/arbolista-dev/cc-user-api

COPY . /go/src/github.com/arbolista-dev/cc-user-api
COPY docker-entrypoint.sh /

RUN chmod u+x docker-entrypoint.sh

ENTRYPOINT ./docker-entrypoint.sh $CC_ENV
