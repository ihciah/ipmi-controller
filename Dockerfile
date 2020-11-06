FROM golang:1.15-alpine AS build

RUN apk add --no-cache git ca-certificates
WORKDIR /ipmi-controller
ADD . / /ipmi-controller/
RUN CGO_ENABLED=0 go build -o /bin/ipmi-controller ./cmd

FROM alpine:latest
MAINTAINER ihciah <ihciah@gmail.com>

RUN apk add --no-cache ca-certificates ipmitool
COPY --from=build /bin/ipmi-controller /bin/ipmi-controller
CMD exec ipmi-controller