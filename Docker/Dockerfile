FROM golang:1.16

RUN apt-get update

COPY ./Social-Network /Social-Network

COPY ./init.sh /

ENTRYPOINT ["/init.sh"]