FROM golang:1.21.5-alpine3.19 as build

COPY . /hnh

WORKDIR /hnh

RUN go build -o ./services/auth/bin/auth ./services/auth/cmd/auth.go

#====================================

FROM alpine:3.19

COPY --from=build /hnh/services/auth/bin/ /bin

RUN mkdir /logs
# WORKDIR /bin

CMD [ "bin/auth" ]
