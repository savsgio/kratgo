###############
### BUILDER ###
###############

FROM golang:1.11-alpine3.9 as builder

RUN apk add git build-base

RUN mkdir -p /usr/src/kratgo
WORKDIR /usr/src/kratgo

ADD . .

RUN make

###############
### RELEASE ###
###############

FROM alpine:3.9

LABEL Author="Sergio Andres Virviescas Santana <developersavsgio@gmail.com>"

COPY --from=builder /usr/src/kratgo/ /kratgo

RUN cd /kratgo \
    && apk add make git \
    && make install \
    && rm -rf /kratgo \
    && apk del make git

# Configuration
COPY ./docker/docker-entrypoint.sh /usr/local/bin/
RUN ln -s /usr/local/bin/docker-entrypoint.sh /entrypoint.sh # backwards compat
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["kratgo"]

EXPOSE 6081 6082
