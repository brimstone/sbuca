FROM scratch

CMD ["server"]

ENTRYPOINT ["/sbuca"]

EXPOSE 8600

ADD sbuca /
