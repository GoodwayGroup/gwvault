FROM alpine:3.13.5

COPY gwvault /usr/local/bin/gwvault
RUN chmod +x /usr/local/bin/gwvault

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/gwvault" ]