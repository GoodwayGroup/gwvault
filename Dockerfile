FROM alpine:3.14.0

COPY gwvault /usr/local/bin/gwvault
RUN chmod +x /usr/local/bin/gwvault

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/gwvault" ]