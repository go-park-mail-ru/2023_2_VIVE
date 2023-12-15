FROM golang:1.21.5-alpine3.19 as build

COPY . /hnh

WORKDIR /hnh

RUN go build -o ./services/searchEngineService/bin/search_engine ./services/searchEngineService/cmd/search_engine.go

#====================================

FROM alpine:3.19

COPY --from=build /hnh/services/searchEngineService/bin/ /bin

WORKDIR /bin

CMD [ "search_engine" ]
