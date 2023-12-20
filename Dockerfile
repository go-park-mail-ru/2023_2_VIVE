FROM golang:1.21.5-alpine3.19 as build

COPY . /hnh

WORKDIR /hnh

RUN apk add make && make build

#====================================

FROM alpine:3.19

WORKDIR /

# COPY /deploy/migrations /migrations
COPY --from=build /hnh/bin/ /bin
COPY --from=build /hnh/assets/ /assets

RUN mkdir /logs


CMD [ "bin/hnh" ]
