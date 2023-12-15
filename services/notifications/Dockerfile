FROM golang:1.21.5-alpine3.19 as build

COPY . /hnh

WORKDIR /hnh

RUN go build -o ./services/notifications/bin/notifications ./services/notifications/cmd/main.go

#====================================

FROM alpine:3.19

COPY --from=build /hnh/services/notifications/bin/ /bin

WORKDIR /bin

CMD [ "notifications" ]
