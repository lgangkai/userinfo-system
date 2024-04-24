FROM golang:1.22.0-alpine
MAINTAINER gangkai
RUN mkdir proto
RUN mkdir userapi
RUN mkdir err
RUN mkdir logger
COPY userapi /userapi
COPY proto /proto
COPY err /err
COPY logger /logger
WORKDIR /proto
RUN go mod download
WORKDIR /err
RUN go mod download
WORKDIR /logger
RUN go mod download
WORKDIR /userapi
RUN go mod download
RUN go build -o userapi .
EXPOSE 8080
ENTRYPOINT sh -c ./userapi -config=conf/userapi.yaml