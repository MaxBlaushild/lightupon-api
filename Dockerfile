FROM golang:1.10

RUN apt-get update -y
RUN apt-get install libmagickwand-dev -y
RUN apt-get install imagemagick -y

RUN mkdir -p /go/src/lightupon-api
WORKDIR /go/src/lightupon-api

ADD . /go/src/lightupon-api

RUN go get -v

EXPOSE 5000

CMD ["go", "run", "main.go"]
