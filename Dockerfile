FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.18

COPY --from=build-stage /app/bin/gwvault /usr/local/bin/gwvault
RUN chmod +x /usr/local/bin/gwvault

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/gwvault" ]
