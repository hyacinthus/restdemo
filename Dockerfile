# build app
FROM golang AS build-env

ADD . /restdemo

WORKDIR /restdemo

RUN go build

# safe image
FROM debian

ENV TZ=Asia/Shanghai

COPY --from=build-env /restdemo/restdemo /usr/bin/restdemo

EXPOSE 1324

CMD ["restdemo"]
