FROM golang:1.11-alpine3.9

LABEL Author="Sergio Andres Virviescas Santana <developersavsgio@gmail.com>"

RUN apk update && apk upgrade
RUN apk add --no-cache git build-base

RUN mkdir -p /usr/src/kratgo
WORKDIR /usr/src/kratgo
ADD . .

RUN make
RUN make install

WORKDIR /

# Clean
RUN rm -rf /usr/src/kratgo
RUN apk del git build-base

# Configuration
COPY ./docker/docker-entrypoint.sh /usr/local/bin/
RUN ln -s /usr/local/bin/docker-entrypoint.sh /entrypoint.sh # backwards compat
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["kratgo"]

EXPOSE 6081 6082
