FROM golang:1.9.0-alpine3.6 as builder

ADD . /go/src/github.com/astapi/pingdom2line
WORKDIR /go/src/github.com/astapi/pingdom2line
RUN apk --update add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

ENV GOPATH /go
RUN go build .

FROM alpine:3.6

COPY --from=builder /go/src/github.com/astapi/pingdom2line/pingdom2line /pingdom2line

RUN apk --update add ca-certificates

ARG line_notify_token
ENV LINE_NOTIFY_TOKEN $line_notify_token

WORKDIR /
RUN chown -R nobody:nogroup /pingdom2line
USER nobody

CMD "./pingdom2line"
