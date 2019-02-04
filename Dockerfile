FROM golang:1.11-alpine3.8

RUN apk add --no-cache git build-base

RUN mkdir -p /usr/src/kratgo
WORKDIR /usr/src/kratgo
ADD . .

RUN make
RUN make install

RUN apk del git build-base

EXPOSE 8000

CMD ["kratgo"]
