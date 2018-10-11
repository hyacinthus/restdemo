# build app
FROM golang AS build-env

ADD . /app

WORKDIR /app

RUN go build

# safe image
FROM debian

ENV TZ=Asia/Shanghai

COPY --from=build-env /app/app /usr/bin/app

EXPOSE 1324

CMD ["app"]
