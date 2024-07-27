FROM alpine:latest
COPY gcredstash .
ENTRYPOINT ["/gcredstash"]
