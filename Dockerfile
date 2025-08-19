##
## Build
##

FROM golang:1.25-buster AS build

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go build -o /air2hugo cmd/air2hugo/main.go
RUN go build -o /dyngo cmd/dyngo/main.go 

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /air2hugo /air2hugo
COPY --from=build /dyngo /dyngo

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/dyngo"]
