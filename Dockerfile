FROM golang:1.18 AS builder

ENV GOPATH /go
ENV APPPATH /repo
COPY . /repo
RUN cd /repo && make

FROM alpine:latest
COPY --from=builder /repo/gcredstash /gcredstash
ENTRYPOINT ["/gcredstash"]
