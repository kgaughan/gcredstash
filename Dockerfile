FROM gcr.io/distroless/static:latest

LABEL org.opencontainers.image.title=gcredstash
LABEL org.opencontainers.image.description="Manages credentials using AWS Key Management Service (KMS) and DynamoDB"
LABEL org.opencontainers.image.vendor="Keith Gaughan"
LABEL org.opencontainers.image.licenses="ASL 2.0"
LABEL org.opencontainers.image.url=https://github.com/kgaughan/gcredstash
LABEL org.opencontainers.image.source=https://github.com/kgaughan/gcredstash
LABEL org.opencontainers.image.documentation=https://kgaughan.github.io/gcredstash/

COPY gcredstash .
ENTRYPOINT ["/gcredstash"]
