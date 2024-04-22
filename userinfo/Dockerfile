FROM golang:1.22.0-alpine
MAINTAINER gangkai
#ENV config conf/userinfo.yaml
RUN mkdir proto
RUN mkdir userinfo
RUN mkdir err
COPY userinfo /userinfo
COPY proto /proto
COPY err /err
WORKDIR /proto
RUN go mod download
WORKDIR /err
RUN go mod download
WORKDIR /userinfo
RUN go mod download
RUN go build -o userinfo .
EXPOSE 8081
#ENTRYPOINT sh -c ./userinfo -config=${config}
ENTRYPOINT sh -c ./userinfo -config=conf/userinfo.yaml