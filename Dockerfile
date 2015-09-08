FROM golang:1.4.1

RUN go get github.com/brimstone/sbuca

VOLUME ["sbuca"]
EXPOSE 8600

ENTRYPOINT ["sbuca"]
CMD ["server"]
