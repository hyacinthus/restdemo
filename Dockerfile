# build
FROM golang:alpine AS build-env

ADD . /art

WORKDIR /art

RUN go build

# app only
FROM alpine

RUN apk add -U tzdata ca-certificates

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

ADD http://static.crandom.com/font/inziu-SC-regular.ttc /usr/share/fonts/

COPY --from=build-env /art/art /usr/bin/art

EXPOSE 1324

CMD ["art"]
