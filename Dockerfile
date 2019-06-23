FROM golang:latest as build-dev

ENV GO111MODULE=on
ENV BUILDPATH=github.com/nsini/blog
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH}/cmd/client && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v
RUN cd /go/src/${BUILDPATH}/cmd/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

COPY --from=build-env /go/bin/server /go/bin/server
COPY --from=build-env /go/bin/client /go/bin/client
 COPY ./static /go/bin/
COPY ./app.yaml /etc/kplcloud/
WORKDIR /go/bin/
CMD ["/go/bin/server", "start", "-p", ":8080", "-c", "/etc/kplcloud/app.yaml"]