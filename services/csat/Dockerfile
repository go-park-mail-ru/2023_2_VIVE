FROM golang:1.21.5-alpine3.19 as build

COPY . /hnh

WORKDIR /hnh

RUN go build -o ./services/csat/bin/csat ./services/csat/cmd/csat.go

#====================================

FROM alpine:3.19

COPY --from=build /hnh/services/csat/bin/ /bin

RUN mkdir /logs
# WORKDIR /bin

CMD [ "bin/csat" ]
