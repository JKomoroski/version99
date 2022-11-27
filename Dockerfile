# Based on https://www.cloudreach.com/en/technical-blog/containerize-this-how-to-build-golang-dockerfiles/
#
# First create main in a rich build environment

FROM golang:alpine as builder
RUN mkdir /src
ADD src /src
WORKDIR /src
RUN go mod init grons.nl/version99 && go build -o version99

# then copy it to a small image.

FROM alpine
MAINTAINER Erik van Oosten, e.vanoosten@grons.nl
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /src/version99 /app/
WORKDIR /app
CMD ["./version99"]
EXPOSE 8080

